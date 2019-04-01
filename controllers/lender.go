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

	lender.User = req.Context().Value("user").(string)

	if err == nil && lender.Create() {
		utils.Response(w, http.StatusOK, types.Response.Ok)
	} else {
		utils.Response(w, http.StatusNotAcceptable, types.Response.InvalidArguments)
	}
}