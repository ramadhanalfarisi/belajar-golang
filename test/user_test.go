package test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
	"tokokocak/app"

	"github.com/google/uuid"
)

var a app.App

func TestMain(m *testing.M) {
	a.Connection("test")

	code := m.Run()
	clearTable()
	os.Exit(code)
}

func clearTable() {
	a.DB.Exec("DELETE FROM users;")
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Mux.ServeHTTP(rr, req)

	return rr
}

func Hashing(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func registerUser(i int) {
	for j := 0; j < i; j++ {
		user_id := uuid.New()
		user_firstname := "user " + strconv.Itoa(j)
		user_lastname := "halo"
		user_email :="user" + strconv.Itoa(j) + "@gmail.com"
		user_address := "address " + strconv.Itoa(j)
		user_password := Hashing("password"+strconv.Itoa(j))
		user_role := "user"
		created_at := time.Now().Format("2006-01-02 15:04:05")
	
		a.DB.Exec("INSERT INTO users VALUES(?, ?, ?, ?, ?, ?, ?, ?, NULL, NULL);",user_id,user_firstname,user_lastname,user_email,user_address,user_password,user_role,created_at)
	}
}

func TestRegisterUser(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{
			"userFirstname" : "Ramadhan",
			"userLastname" : "Alfarisi",
			"userEmail" : "ramadhan@gmail.com",
			"userAddress" : "Banyuwangi",
			"userPassword" : "konohagakure",
			"userRepassword" : "konohagakure",
			"userRole" : "user"
		}`)
	req, _ := http.NewRequest("POST", "/v1/register", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestLoginUser(t *testing.T) {
	clearTable()
	registerUser(1)
	var jsonStr = []byte(`{
			"userEmail" : "user1@gmail.com",
			"userPassword" : "password1"
		}`)
	req, _ := http.NewRequest("POST", "/v1/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
