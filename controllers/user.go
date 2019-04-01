package controllers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"p2p-lending/models"
	"p2p-lending/types"
	"p2p-lending/utils"
)

func Login (w http.ResponseWriter, r *http.Request) {

	parser := &models.User{}
	err := json.NewDecoder(r.Body).Decode(parser) //decode the request body into struct and failed if any error occur
	if err != nil {
		utils.Response(w, http.StatusNotAcceptable, types.Response.EmailOrPasswordInvalid)
		return
	}

	email := parser.Email
	password := parser.Password

	user := &models.User{}
	err = models.GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil {
		//if err == gorm.ErrRecordNotFound {
		//	utils.Response(w, http.StatusNotFound, types.Response.NotFound)
		//	return
		//}
		utils.Response(w, http.StatusNotFound, types.Response.EmailOrPasswordInvalid)
		return
	}

	//GetDB().Model(&account).Association("Environments").Find(&account.Environments)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		utils.Response(w, http.StatusNotFound, types.Response.EmailOrPasswordInvalid)
		return
	}

	//Create JWT token
	tk := &models.Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))

	var result = struct{
		Token string `json:"token"`
	}{tokenString}

	utils.ResponseJson(w, http.StatusNotAcceptable, result)
}

func GetUser(w http.ResponseWriter, req *http.Request) {
	user := models.GetUserById(req.Context().Value("user").(string))

	if user.ID != "" {
		utils.ResponseJson(w, http.StatusOK, user)
	} else {
		var s struct{}
		utils.ResponseJson(w, http.StatusOK, s)
	}
}

func GetUserById(w http.ResponseWriter, req *http.Request) {
	user := models.GetUserById(mux.Vars(req)["id"])

	if user.ID != "" {
		var result = struct{
			Score int
			Salary float32
		}{Score: user.Score, Salary: user.Salary}

		utils.ResponseJson(w, http.StatusOK, result)
	} else {
		var s struct{}
		utils.ResponseJson(w, http.StatusOK, s)
	}
}

func AddUser(w http.ResponseWriter, req *http.Request) {
	user := models.User{}
	err := json.NewDecoder(req.Body).Decode(&user)

	if err == nil && user.Create() {
		utils.Response(w, http.StatusOK, types.Response.Ok)
	} else {
		utils.Response(w, http.StatusNotAcceptable, types.Response.InvalidArguments)
	}
}

func PayLending(w http.ResponseWriter, req *http.Request) {
	user := models.GetUserById(req.Context().Value("user").(string))
	payment := mux.Vars(req)["id"]

	utils.ResponseJson(w, http.StatusOK, user.Pay(payment))
}