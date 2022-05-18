package controllers

import (
	"belajar-golang/helpers"
	"belajar-golang/models"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func Hashing(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func Login(w http.ResponseWriter, r *http.Request){
	db, err := helpers.Connection()
	if err != nil {
		log.Fatal(err)
	}
	var userLogin models.UserLogin
	err_json := json.NewDecoder(r.Body).Decode(&userLogin)
	if err_json != nil{
		log.Fatal(err)
	}
	msg, isvalid := helpers.Validate(userLogin)
	if isvalid == false{
		response := helpers.InvalidResponse(400, msg)
		json, err := json.Marshal(response)
		if err != nil{
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
	}else{
		userLogin.UserPassword = Hashing(userLogin.UserPassword)
		login, err := userLogin.LoginUser(db)
		if err != nil{
			log.Fatal(err)
		}
		jwt := helpers.GenerateJWT(login)
		response := map[string]interface{}{
			"code" : 200,
			"token" : jwt,
			"message" : "Login succcessfully",
		}
		json, err := json.Marshal(response)
		if err != nil{
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func Register(w http.ResponseWriter, r *http.Request){
	db, err := helpers.Connection()
	if err != nil {
		log.Fatal(err)
	}
	var userRegister models.UserRegister
	userRegister.UserId = uuid.New()
	err_json := json.NewDecoder(r.Body).Decode(&userRegister)
	if err_json != nil{
		log.Fatal(err)
	}
	msg, isvalid := helpers.Validate(userRegister)
	if isvalid == false{
		response := helpers.InvalidResponse(400, msg)
		json, err := json.Marshal(response)
		if err != nil{
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
	}else{
		var user models.User
		user.UserId = userRegister.UserId
		user.UserFirstname = userRegister.UserFirstname
		user.UserLastname = userRegister.UserLastname
		user.UserAddress = userRegister.UserAddress
		user.UserEmail = userRegister.UserEmail
		user.UserPassword = Hashing(userRegister.UserPassword)
		user.UserRole = userRegister.UserRole
		err := user.RegisterUser(db)
		if err != nil{
			log.Fatal(err)
		}
		response := helpers.SuccessResponse(200,"Register successfully",nil)
		json, err := json.Marshal(response)
		if err != nil{
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}