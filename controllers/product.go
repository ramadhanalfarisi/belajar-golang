package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"tokokocak/helpers"
	"tokokocak/middlewares"
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
	}

	userId := helpers.GetUserId(r)
	searchRedis := helpers.SearchRedisValue("belajar:product:" + userId.String())
	if searchRedis != nil {
		product := helpers.GetRedisValue("belajar:product:" + userId.String())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(product))
	} else {
		var products models.Product
		var get_product_response GetProductResponse
		products.UserId = userId
		meta_param := r.Context().Value("metaParam").(middlewares.MetaParam)

		int_limit := meta_param.Limit
		int_offset := meta_param.Offset

		var func_get_product = func(ch chan []models.GetProduct) {
			get_product, err := products.SelectAllProducts(db, []string{"product_id", "product_name", "product_desc", "product_price", "product_image", "created_at", "updated_at"}, int_limit, int_offset)
			if err != nil {
				log.Println(err)
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

		response := helpers.SuccessResponse(200, "Get products successfully", get_product_response.Item, get_product_response.Meta)
		json, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
		}
		helpers.SetRedisValue("belajar:product:" + userId.String(), string(json))
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func SelectOneProduct(w http.ResponseWriter, r *http.Request) {
	db, err := helpers.Connection()

	if err != nil {
		log.Println(err)
	}

	params := mux.Vars(r)
	id := params["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Println(err)
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
		if err != nil{
			log.Println(err)
		}
		response := helpers.SuccessResponse(200,"Get product successfully",get_product,nil)
		json, err := json.Marshal(response)
		if err != nil{
			log.Println(err)
		}
		helpers.SetRedisValue("belajar:product:" + userId.String() + ":" + id, string(json))
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}
