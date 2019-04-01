package controllers

import (
	"net/http"
	"p2p-lending/models"
	"p2p-lending/utils"
)

func GetStatements(w http.ResponseWriter, req *http.Request) {
	userID := req.Context().Value("user").(string)
	statements := models.GetStatementsByUser(userID)
	utils.ResponseJson(w, http.StatusOK, statements)
}