package middlewares

import (
	"context"
	"log"
	"net/http"
	"strconv"
)

type MetaParam struct {
	Limit  int64
	Page   int64
	Offset int64
}

func ApiMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			var int_limit int64
			var int_page int64
			var int_offset int64

			page, ok := r.URL.Query()["page"]
			limit, ok := r.URL.Query()["limit"]
			offset, ok := r.URL.Query()["offset"]

			climit, err1 := strconv.Atoi(limit[0])
			cpage, err2 := strconv.Atoi(page[0])
			coffset, err3 := strconv.Atoi(offset[0])

			if err1 != nil {
				log.Println(err1)
			}

			if err2 != nil {
				log.Println(err2)
			}

			if err3 != nil {
				log.Println(err3)
			}

			if !ok || len(limit[0]) < 1 {
				int_limit = 10
			} else {
				int_limit = int64(climit)
			}

			if !ok || len(page[0]) < 1 {
				int_page = 1
			} else {
				int_page = int64(cpage)
			}

			if !ok || len(offset[0]) < 1 {
				int_offset = 0
			} else {
				int_offset = int64(coffset)
			}
			var meta_param MetaParam
			meta_param.Page = int_page
			meta_param.Limit = int_limit
			meta_param.Offset = int_offset
			ctx := context.WithValue(context.Background(), "metaParam", meta_param)
			r = r.WithContext(ctx)
		}

		w.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}
