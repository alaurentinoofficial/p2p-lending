package controllers

import (
	"encoding/json"
	"net/http"
	"p2p-lending/models"
	"p2p-lending/types"
	"p2p-lending/utils"
)

func AddUser(w http.ResponseWriter, req *http.Request) {
	user := models.User{}
	err := json.NewDecoder(req.Body).Decode(&user)

	if err == nil && user.Create() {
		utils.Response(w, http.StatusOK, types.Response.Ok)
	} else {
		utils.Response(w, http.StatusNotAcceptable, types.Response.InvalidArguments)
	}
}