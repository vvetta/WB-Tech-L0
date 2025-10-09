package http

import (
	"net/http"
)

func NewRouter(d Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/order/", d.handleGetOrderByID)

	return withCORS(mux)
}
