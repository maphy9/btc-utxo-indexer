package requests

import (
	"encoding/json"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewRegisterRequest(r *http.Request) (RegisterRequest, error) {
	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}