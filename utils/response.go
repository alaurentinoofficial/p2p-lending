package utils

import (
	"encoding/json"
	"net/http"
)

type ResponseType struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var ResponseMap = make(map[int]string)

func init() {
	ResponseMap[0] = "Successfully"
	ResponseMap[1] = "Unauthorized"
	ResponseMap[2] = "Invalid Arguments"
	ResponseMap[3] = "Not Found"
	ResponseMap[4] = "Already Exists"
}

func Response(w http.ResponseWriter, status int, code int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/json")
	_ = json.NewEncoder(w).Encode(ResponseType{Message: ResponseMap[code], Code: code})
}

func ResponseJson(w http.ResponseWriter, status int, obj interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/json")
	_ = json.NewEncoder(w).Encode(obj)
}