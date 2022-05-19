package middlewares

import (
	"belajar-golang/helpers"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			authorization := r.Header.Get("Authorization")
			if !strings.Contains(authorization, "Bearer") {
				response := helpers.FailedResponse(401, "Token must be Bearer type")
				json, err := json.Marshal(response)
				if err != nil {
					log.Fatal(err)
				}
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(json)
			} else {
				tokenString := strings.Replace(authorization, "Bearer ", "", -1)
				decodeToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
					if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Signing Method invalid")
					} else if method != helpers.JWT_SIGNING_METHOD {
						return nil, fmt.Errorf("Signing Method invalid")
					}

					err_claims := t.Claims.Valid()
					if err_claims != nil {
						return nil, err_claims
					}

					return helpers.JWT_SIGNATURE_KEY, nil
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				claims, ok := decodeToken.Claims.(jwt.MapClaims)

				if !ok || !decodeToken.Valid {
					http.Error(w, "Not Valid", http.StatusBadGateway)
				}

				ctx := context.WithValue(context.Background(), "userDetail", claims)
				r = r.WithContext(ctx)
				handler.ServeHTTP(w, r)
			}
		} else {
			response := helpers.FailedResponse(401, "Authorization header is required")
			json, err := json.Marshal(response)
			if err != nil {
				log.Fatal(err)
			}
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(json)
		}
	})
}
