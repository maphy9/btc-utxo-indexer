package apierrors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/jsonapi"
)

func NewApiError(status int, detail string) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title: http.StatusText(status),
		Detail: detail,
		Status: fmt.Sprintf("%d", status),
		Code: strings.ReplaceAll(strings.ToLower(detail), " ", "_"),
	}
}

func BadRequest() *jsonapi.ErrorObject {
	return NewApiError(http.StatusBadRequest, "Bad request")
}