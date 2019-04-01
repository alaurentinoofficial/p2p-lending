package controllers

import (
	"encoding/json"
	"net/http"
	"p2p-lending/models"
	"p2p-lending/types"
	"p2p-lending/utils"
)

func GetLenders(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value("user").(string)
	lenders := models.GetLendersByUser(userID)

	utils.ResponseJson(w, http.StatusOK, lenders)
}

func AddLender(w http.ResponseWriter, req *http.Request) {
	lender := models.Lender{}
	err := json.NewDecoder(req.Body).Decode(&lender)

	if err == nil {
		utils.Response(w, http.StatusOK, lender.Create())
	} else {
		utils.Response(w, http.StatusNotAcceptable, types.Response.InvalidArguments)
	}
}