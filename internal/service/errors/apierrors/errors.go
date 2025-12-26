package apierrors

import (
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
)

func NewApiError(status int, detail string, code string) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title: http.StatusText(status),
		Detail: detail,
		Status: fmt.Sprintf("%d", status),
		Code: code,
	}
}

func BadRequest() *jsonapi.ErrorObject {
	return NewApiError(http.StatusBadRequest, "Bad request", "bad_request	")
}