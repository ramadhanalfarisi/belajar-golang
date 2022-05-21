package helpers

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func GetUserId(r *http.Request) uuid.UUID {
	userDetail := r.Context().Value("userDetail").(jwt.MapClaims)
	userId := userDetail["userId"].(string)
	uuid, err := uuid.Parse(userId)
	if err != nil {
		log.Println(err)
	}
	return uuid
}
