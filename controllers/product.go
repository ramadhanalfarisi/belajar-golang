package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
	"tokokocak/helpers"
	"tokokocak/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type GetProductResponse struct {
	Item []models.GetProduct
	Meta models.Pagination
}

func (controller *Controller) SelectAllProducts(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(2)
	db := controller.DB
	userId := helpers.GetUserId(r)
	meta_param, _ := models.GetMetaParam(r)
	int_limit := meta_param.Limit
	int_offset := meta_param.Offset
	string_limit := strconv.Itoa(int(int_limit))
	string_offset := strconv.Itoa(int(int_offset))
	searchRedis := helpers.SearchRedisValue("tokokocak:product:" + userId.String() + ":" + string_limit + ":" + string_offset)
	if len(searchRedis) > 0 {
		product := helpers.GetRedisValue("tokokocak:product:" + userId.String() + ":" + string_limit + ":" + string_offset)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(product))
	} else {
		var products models.Product
		var get_product_response GetProductResponse
		products.UserId = userId

		var func_get_product = func(ch chan []models.GetProduct) {
			get_product, err := products.SelectAllProducts(db, []string{"product_id", "product_name", "product_desc", "product_price", "product_image", "created_at", "updated_at"}, int_limit, int_offset)
			if err != nil {
				helpers.Error(err)
				return
			}
			ch <- get_product
		}

		var func_get_meta = func(ch chan models.Pagination) {
			get_product_num := products.SelectRowProducts(db, []string{"product_id"})
			var pagination models.Pagination
			pagination.Total = get_product_num
			result_pagination := pagination.CreatePagination(r, meta_param)
			ch <- result_pagination
		}

		var chProd = make(chan []models.GetProduct)
		go func_get_product(chProd)

		var chMeta = make(chan models.Pagination)
		go func_get_meta(chMeta)

		for i := 0; i < 2; i++ {
			select {
			case prod := <-chProd:
				get_product_response.Item = prod
			case meta := <-chMeta:
				get_product_response.Meta = meta
			}
		}

		response := helpers.SuccessResponse(http.StatusOK, "Get products successfully", get_product_response.Item, get_product_response.Meta)
		json, err := json.Marshal(response)
		if err != nil {
			helpers.Error(err)
			return
		}
		helpers.SetRedisValue("tokokocak:product:"+userId.String()+":"+string_limit+":"+string_offset, string(json))
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (controller *Controller) SelectOneProduct(w http.ResponseWriter, r *http.Request) {
	db := controller.DB
	params := mux.Vars(r)
	id := params["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		helpers.Error(err)
		return
	}

	var product models.Product
	product.ProductId = uuid

	userId := helpers.GetUserId(r)
	searchRedis := helpers.SearchRedisValue("tokokocak:product:" + userId.String() + ":" + id)
	if len(searchRedis) > 0 {
		product := helpers.GetRedisValue("tokokocak:product:" + userId.String() + ":" + id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(product))
	} else {
		product.UserId = userId
		get_product, err := product.SelectOneProduct(db, []string{"product_id", "product_name", "product_desc", "product_price", "product_image", "created_at", "updated_at"})
		if err != nil {
			helpers.Error(err)
			return
		}
		response := helpers.SuccessResponse(http.StatusOK, "Get product successfully", get_product, nil)
		json, err := json.Marshal(response)
		if err != nil {
			helpers.Error(err)
			return
		}
		helpers.SetRedisValue("tokokocak:product:"+userId.String()+":"+id, string(json))
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (controller *Controller) InsertProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	err2 := json.NewDecoder(r.Body).Decode(&products)
	if err2 != nil {
		helpers.Error(err2)
		return
	}

	userId := helpers.GetUserId(r)
	var validateall []string
	for i := range products {
		products[i].ProductId = uuid.New()
		products[i].UserId = userId
		products[i].CreatedAt = time.Now().Format("2006-01-02 15:04:05")
		products[i].UpdatedAt = nil
		validate, _ := helpers.Validate(products[i])
		validateall = append(validateall, validate...)
	}

	if validateall == nil {
		chanProduct := convertChanProduct(products)
		chanInsert1 := controller.insertProduct(chanProduct)
		chanInsert2 := controller.insertProduct(chanProduct)
		mergeChanOut := mergeChanOut(chanInsert1, chanInsert2)
		for message := range mergeChanOut {
			fmt.Println(message)
		}
		response := helpers.SuccessResponse(200, "Insert product succesfully", nil, nil)
		json, err := json.Marshal(response)
		if err != nil {
			helpers.Error(err)
		}
		getredis := helpers.SearchRedisValue("tokokocak:product:" + userId.String())
		if len(getredis) > 0 {
			helpers.DeleteRedisValue(getredis)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	} else {
		response := helpers.InvalidResponse(http.StatusBadRequest, validateall)
		json, err := json.Marshal(response)
		if err != nil {
			helpers.Error(err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
	}
}

func convertChanProduct(products []models.Product) <-chan models.Product {
	chanOut := make(chan models.Product)

	go func() {
		for _, product := range products {
			chanOut <- product
		}
		close(chanOut)
	}()

	return chanOut
}

func (controller *Controller) insertProduct(chanProd <-chan models.Product) <-chan string {
	chanOut := make(chan string)

	go func() {
		for prod := range chanProd {
			prod.InsertProduct(controller.DB)
			chanOut <- "Insert data successfully"
		}
		close(chanOut)
	}()
	return chanOut
}

func mergeChanOut(chanInserts ...<-chan string) <-chan string {
	wg := new(sync.WaitGroup)
	chanOut := make(chan string)

	wg.Add(len(chanInserts))
	for _, chanInsert := range chanInserts {
		go func(chanIns <-chan string) {
			for chanIn := range chanIns {
				chanOut <- chanIn
			}
			wg.Done()
		}(chanInsert)
	}

	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}

func (controller *Controller) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := controller.DB
	userId := helpers.GetUserId(r)
	params := mux.Vars(r)
	id := params["id"]
	uuid, _ := uuid.Parse(id)

	var product models.Product
	err2 := json.NewDecoder(r.Body).Decode(&product)
	if err2 != nil {
		helpers.Error(err2)
	}
	product.ProductId = uuid
	product.UserId = userId

	validate, _ := helpers.Validate(product)

	if validate == nil {
		now := time.Now().Format("2006-01-02 15:04:05")
		data := models.Product{
			ProductName:  product.ProductName,
			ProductDesc:  product.ProductDesc,
			ProductPrice: product.ProductPrice,
			ProductImage: product.ProductImage,
			UpdatedAt:    &now,
		}
		res_products, err := product.UpdateProduct(data, db)
		if err != nil {
			helpers.Error(err)
		}
		response := helpers.SuccessResponse(200, "Update product successfully", res_products, nil)
		json, _ := json.Marshal(response)
		getredis := helpers.SearchRedisValue("tokokocak:product:" + userId.String())
		if len(getredis) > 0 {
			helpers.DeleteRedisValue(getredis)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	} else {
		response := helpers.InvalidResponse(500, validate)
		json, _ := json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
	}
}

func (controller *Controller) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	connect := controller.DB
	userId := helpers.GetUserId(r)
	params := mux.Vars(r)
	id := params["id"]
	uuid, _ := uuid.Parse(id)

	var product models.Product
	err2 := json.NewDecoder(r.Body).Decode(&product)
	if err2 != nil {
		helpers.Error(err2)
	}
	product.ProductId = uuid
	product.UserId = userId

	err_delete := product.DeleteProduct(connect)
	if err_delete != nil {
		helpers.Error(err_delete)
	}
	response := helpers.SuccessResponse(200, "Delete product successfully", nil, nil)
	json, _ := json.Marshal(response)
	getredis := helpers.SearchRedisValue("tokokocak:product:" + userId.String())
	if len(getredis) > 0 {
		helpers.DeleteRedisValue(getredis)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
