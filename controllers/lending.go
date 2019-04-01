package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"p2p-lending/models"
	"p2p-lending/types"
	"p2p-lending/utils"
)

func getLendings(w http.ResponseWriter, req *http.Request) {
	lendings := models.GetLendings()
	utils.ResponseJson(w, http.StatusOK, lendings)
}

func getLendingById(w http.ResponseWriter, req *http.Request) {
	lending := models.GetLendingById(mux.Vars(req)["id"])

	if lending.ID != "" {
		utils.ResponseJson(w, http.StatusOK, lending)
	} else {
		utils.Response(w, http.StatusOK, types.Response.NotFound)
	}
}

func AddLending(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value("user").(string)

	lending := models.Lending{}
	err := json.NewDecoder(req.Body).Decode(&lending)
	lending.Taker = userID

	if err != nil && lending.Create() {
		utils.Response(w, http.StatusOK, types.Response.Ok)
	}
}