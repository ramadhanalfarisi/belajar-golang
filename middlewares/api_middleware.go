package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"tokokocak/helpers"
	"tokokocak/models"
)

type MetaParam struct {
	Limit  int64
	Page   int64
	Offset int64
}

func ApiMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			meta_param, err := models.GetMetaParam(r)
			if err != nil{
				response := helpers.FailedResponse(http.StatusBadRequest,err.Error())
				json,err := json.Marshal(response)
				if err != nil{
					log.Println(err)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				w.Write(json)
			}
			ctx := context.WithValue(context.Background(), "metaParam", meta_param)
			r = r.WithContext(ctx)
		}

		w.Header().Add("Content-Type", "application/json")
		handler.ServeHTTP(w, r)
	})
}
