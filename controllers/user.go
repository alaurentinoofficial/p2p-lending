package controllers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
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
		utils.Response(w, http.StatusNotAcceptable, types.Response.InvalidArguments)
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
		utils.Response(w, http.StatusNotFound, types.Response.NotFound)
		return
	}

	//GetDB().Model(&account).Association("Environments").Find(&account.Environments)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		utils.Response(w, http.StatusNotFound, types.Response.NotFound)
		return
	}

	//Create JWT token
	tk := &models.Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))

	utils.ResponseJson(w, http.StatusNotAcceptable, models.TokenResponse{Token: tokenString})
}

func GetUser(w http.ResponseWriter, req *http.Request) {
	user := models.GetUserById(req.Context().Value("user").(string))

	utils.ResponseJson(w, http.StatusOK, user)
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