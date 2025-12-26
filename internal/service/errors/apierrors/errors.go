package apierrors

import (
	"net/http"

	"github.com/google/jsonapi"
)

func BadRequest() *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(400),
		Detail: "Bad Request",
		Status: "400",
		Code:   "bad_request",
	}
}

func InternalServerError() *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(500),
		Detail: "Internal Server Error",
		Status: "500",
		Code:   "internal_server_error",
	}
}

func NotFound() *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(404),
		Detail: "Not Found",
		Status: "404",
		Code:   "no_mapping_found",
	}
}