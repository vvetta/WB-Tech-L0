package http

import (
	"errors"
	"strings"
	"net/http"

	"WB-Tech-L0/internal/domain"
)

func (d Deps) handleGetOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")	
		return
	}

	orderUID := strings.TrimPrefix(r.URL.Path, "/order/")
	if orderUID == "" || orderUID == r.URL.Path {
		writeError(w, http.StatusBadRequest, "missing order_uid")	
		return
	}

	order, err := d.OrderSvc.GetByID(r.Context(), orderUID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "order not found")
			return
		}
		d.Logger.Error("http.handleGetOrderByID: get order failed", "order_uid", orderUID, "err", err)
		
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	
	writeJSON(w, http.StatusOK, toHttpOrder(order))
}
