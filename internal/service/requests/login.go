package requests

import (
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewLoginRequest(r *http.Request) (LoginRequest, error) {
	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}
