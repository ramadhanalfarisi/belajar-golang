package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"
	"tokokocak/helpers"
	"tokokocak/models"

	"github.com/google/uuid"
)

func Hashing(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}


func (controller *Controller) Login(w http.ResponseWriter, r *http.Request){
	db := controller.DB
	
	var userLogin models.UserLogin
	err_json := json.NewDecoder(r.Body).Decode(&userLogin)
	if err_json != nil{
		helpers.Error(err_json)
		return
	}
	msg, isvalid := helpers.Validate(userLogin)
	if !isvalid{
		response := helpers.InvalidResponse(400, msg)
		json, err := json.Marshal(response)
		if err != nil{
			helpers.Error(err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
	}else{
		userLogin.UserPassword = Hashing(userLogin.UserPassword)
		login, err := userLogin.LoginUser(db)
		if err != nil{
			helpers.Error(err)
			return
		}
		jwt := helpers.GenerateJWT(login)
		response := map[string]interface{}{
			"code" : 200,
			"token" : jwt,
			"message" : "Login succcessfully",
		}
		json, err := json.Marshal(response)
		if err != nil{
			helpers.Error(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (controller *Controller) Register(w http.ResponseWriter, r *http.Request){
	db := controller.DB
	
	var userRegister models.UserRegister
	err_json := json.NewDecoder(r.Body).Decode(&userRegister)
	if err_json != nil{
		helpers.Error(err_json)
		return
	}
	msg, isvalid := helpers.Validate(userRegister)
	if !isvalid{
		response := helpers.InvalidResponse(400, msg)
		json, err := json.Marshal(response)
		if err != nil{
			helpers.Error(err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(json)
	}else{
		var user models.User
		user.UserId = uuid.New()
		user.UserFirstname = userRegister.UserFirstname
		user.UserLastname = userRegister.UserLastname
		user.UserAddress = userRegister.UserAddress
		user.UserEmail = userRegister.UserEmail
		user.UserPassword = Hashing(userRegister.UserPassword)
		user.UserRole = userRegister.UserRole
		user.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
		user.UpdatedAt = nil
		err := user.RegisterUser(db)
		if err != nil{
			helpers.Error(err)
			return
		}
		response := helpers.SuccessResponse(200,"Register successfully",nil,nil)
		json, err := json.Marshal(response)
		if err != nil{
			helpers.Error(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}