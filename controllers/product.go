package controllers

import (
	"belajar_golang/helpers"
	"belajar_golang/middlewares"
	"belajar_golang/models"
	"encoding/json"
	"log"
	"net/http"
	"runtime"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type GetProductResponse struct {
	Item []models.GetProduct
	Meta models.Pagination
}

func SelectAllProducts(w http.ResponseWriter, r *http.Request) {
	runtime.GOMAXPROCS(2)
	db, err := helpers.Connection()
	if err != nil {
		log.Fatal(err)
	}

	userDetail := r.Context().Value("userDetail").(jwt.MapClaims)
	userId := userDetail["userId"].(string)
	searchRedis := helpers.SearchRedisValue("belajar:product:" + userId)
	if searchRedis != nil {
		product := helpers.GetRedisValue("belajar:product:" + userId)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(product))
	} else {
		var products models.Product
		var get_product_response GetProductResponse
		userid, err := uuid.Parse(userId)
		if err != nil {
			log.Fatal(err)
		}
		products.UserId = userid
		meta_param := r.Context().Value("metaParam").(middlewares.MetaParam)

		int_limit := meta_param.Limit
		int_offset := meta_param.Offset

		var func_get_product = func(ch chan []models.GetProduct) {
			get_product, err := products.SelectAllProducts(db, []string{"product_id", "product_name", "product_desc", "product_price", "product_image", "created_at", "updated_at"}, int_limit, int_offset)
			if err != nil {
				log.Fatal(err)
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
		json,err := json.Marshal(response)
		if err != nil{
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}
