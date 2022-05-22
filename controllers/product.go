package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"tokokocak/helpers"
	"tokokocak/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type GetProductResponse struct {
	Item []models.GetProduct
	Meta models.Pagination
}

func SelectAllProducts(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(2)
	db, err := helpers.Connection()
	if err != nil {
		log.Println(err)
		return
	}

	userId := helpers.GetUserId(r)
	meta_param := r.Context().Value("metaParam").(models.MetaParam)

	int_limit := meta_param.Limit
	int_offset := meta_param.Offset
	string_limit := strconv.Itoa(int(int_limit))
	string_offset := strconv.Itoa(int(int_offset))
	searchRedis := helpers.SearchRedisValue("belajar:product:" + userId.String() + ":" + string_limit + ":" + string_offset)
	if searchRedis != nil {
		product := helpers.GetRedisValue("belajar:product:" + userId.String() + ":" + string_limit + ":" + string_offset)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(product))
	} else {
		var products models.Product
		var get_product_response GetProductResponse
		products.UserId = userId

		var func_get_product = func(ch chan []models.GetProduct) {
			get_product, err := products.SelectAllProducts(db, []string{"product_id", "product_name", "product_desc", "product_price", "product_image", "created_at", "updated_at"}, int_limit, int_offset)
			if err != nil {
				log.Println(err)
				return
			}
			ch <- get_product
		}

		var func_get_meta = func(ch chan models.Pagination) {
			get_product_num := products.SelectRowProducts(db, []string{"product_id"})
			var pagination models.Pagination
			pagination.Total = get_product_num
			result_pagination := pagination.CreatePagination(r)
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
			log.Println(err)
			return
		}
		helpers.SetRedisValue("belajar:product:"+userId.String()+":"+string_limit+":"+string_offset, string(json))
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func SelectOneProduct(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.Connection()

	if err != nil {
		log.Println(err)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Println(err)
		return
	}

	var product models.Product
	product.ProductId = uuid

	userId := helpers.GetUserId(r)
	searchRedis := helpers.SearchRedisValue("belajar:product:" + userId.String() + ":" + id)
	if searchRedis != nil {
		product := helpers.GetRedisValue("belajar:product:" + userId.String() + ":" + id)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(product))
	} else {
		product.UserId = userId
		get_product, err := product.SelectOneProduct(db, []string{"product_id", "product_name", "product_desc", "product_price", "product_image", "created_at", "updated_at"})
		if err != nil {
			log.Println(err)
			return
		}
		response := helpers.SuccessResponse(http.StatusOK, "Get product successfully", get_product, nil)
		json, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			return
		}
		helpers.SetRedisValue("belajar:product:"+userId.String()+":"+id, string(json))
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func InsertProducts(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.Connection()

	if err != nil {
		log.Println(err)
	}

	var products []models.Product

	err2 := json.NewDecoder(r.Body).Decode(&products)
	if err2 != nil {
		log.Println(err2)
		return
	}

	userId := helpers.GetUserId(r)
	var isvalidall bool = true
	var validateall []string
	for i, product := range products {
		products[i].ProductId = uuid.New()
		products[i].UserId = userId
		validate, isvalid := helpers.Validate(product)
		if isvalidall == true {
			isvalidall = isvalid
		}
		validateall = append(validateall, validate...)
	}

	if isvalidall {
		insert_func := func(ch chan<- string) {
			size := 5
			var j int
			for i := 0; i < len(products); i += size {
				j += size
				if j > len(products) {
					j = len(products)
				}
				res_products := products[i:j]
				insert_err := models.InsertProduct(res_products, db)
				if insert_err != nil {
					log.Println(insert_err)
					return
				}
				ch <- fmt.Sprintf("Products have inserted %d:%d", i, j)
			}
			close(ch)
		}
		print_msg := func (ch <-chan string)  {
			for message := range ch {
				fmt.Println(message)
			}
		}
		var messages = make(chan string)
		go insert_func(messages)
		
		print_msg(messages)
		response := helpers.SuccessResponse(200, "Insert product succesfully", nil,nil)
		json, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		getredis := helpers.SearchRedisValue("belajar:product:"+userId.String())
		if len(getredis) > 0 {
			helpers.DeleteRedisValue(getredis)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	} else {
		response := helpers.InvalidResponse(http.StatusBadRequest, validateall)
		json, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
	}
}
