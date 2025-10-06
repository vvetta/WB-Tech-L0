package http

import (
	"net/http"
)

func NewRouter(d Deps) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/order/", d.handleGetOrderByID)

	return mux	
}
