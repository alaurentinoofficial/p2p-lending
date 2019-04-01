package controllers

import (
	"encoding/json"
	"net/http"
	"p2p-lending/models"
	"p2p-lending/types"
	"p2p-lending/utils"
)

func AddLending(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value("user").(string)

	lending := models.Lending{}
	err := json.NewDecoder(req.Body).Decode(&lending)
	lending.Taker = userID

	if err != nil && lending.Create() {
		utils.Response(w, http.StatusOK, types.Response.Ok)
	}
}