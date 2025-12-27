package requests

import (
	"encoding/json"
	"net/http"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewRefreshRequest(r *http.Request) (RefreshRequest, error) {
	var request RefreshRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}
