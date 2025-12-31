package util

import (
	"io"
	"net/http"
	"strconv"
)

func ParseString(res *http.Response) (string, error) {
	body, err := io.ReadAll(res.Body)
	return string(body), err
}

func ParseInt64(res *http.Response) (int64, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(string(body), 10, 64)
}
