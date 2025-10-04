package http

import (
	"net/http"
)

func NewRouter(d Deps) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/orders/", d.handleGetOrderByID)

	return mux	
}
