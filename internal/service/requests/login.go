package requests

import "net/http"

type LoginRequest struct {

}

func NewLoginRequest(r *http.Request) (LoginRequest, error) {
	return LoginRequest{}, nil
}