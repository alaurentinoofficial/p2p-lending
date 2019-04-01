package utils

import (
	"encoding/json"
	"net/http"
)

type ResponseType struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var ResponseMap map[int]string

func init() {
	ResponseMap[0] = "Successfully"
	ResponseMap[1] = "Unauthorized"
	ResponseMap[2] = "InvalidArguments"
	ResponseMap[3] = "NotFound"
}

func Response(w http.ResponseWriter, status int, code int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/json")
	_ = json.NewEncoder(w).Encode(ResponseType{Message: ResponseMap[code], Code: code})
}