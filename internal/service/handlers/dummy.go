package handlers

import "net/http"

func Dummy(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("IT IS WORKING"))
}
