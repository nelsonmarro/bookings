package tests

import (
	"net/http"
)

type MyHandler struct{}

func (mh *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Your handler logic here
}
