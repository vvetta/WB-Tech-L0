package http_test

import (
	httpapi "WB-Tech-L0/internal/adapters/http"
	"WB-Tech-L0/internal/domain"
	"WB-Tech-L0/internal/usecase/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type noopLogger struct {}
func (noopLogger) Info(string, ...any) {}
func (noopLogger) Error(string, ...any) {}
func (noopLogger) Debug(string, ...any) {}

func TestGetOrder_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mocks.NewMockOrderReader(ctrl)
	
	order := &domain.Order{OrderUID: "id-1"}
	reader.EXPECT().GetByID(gomock.Any(), "id-1").Return(order, nil)

	mux := httpapi.NewRouter(httpapi.Deps{OrderSvc: reader, Logger: noopLogger{}})
	req := httptest.NewRequest(http.MethodGet, "/order/id-1", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"order_uid":"id-1"`)
}

func TestGetOrder_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mocks.NewMockOrderReader(ctrl)
	
	reader.EXPECT().GetByID(gomock.Any(), "notFoundID").Return(nil, domain.ErrNotFound)

	mux := httpapi.NewRouter(httpapi.Deps{OrderSvc: reader, Logger: noopLogger{}})
	req := httptest.NewRequest(http.MethodGet, "/order/notFoundID", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetOrder_MethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mocks.NewMockOrderReader(ctrl)

	mux := httpapi.NewRouter(httpapi.Deps{OrderSvc: reader, Logger: noopLogger{}})
	req := httptest.NewRequest(http.MethodPost, "/order/", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestGetOrder_EmptyOrderUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mocks.NewMockOrderReader(ctrl)

	mux := httpapi.NewRouter(httpapi.Deps{OrderSvc: reader, Logger: noopLogger{}})
	req := httptest.NewRequest(http.MethodGet, "/order/", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOrder_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mocks.NewMockOrderReader(ctrl)
	
	dbErr := errors.New("db error")

	reader.EXPECT().GetByID(gomock.Any(), "id-1").Return(nil, dbErr)

	mux := httpapi.NewRouter(httpapi.Deps{OrderSvc: reader, Logger: noopLogger{}})
	req := httptest.NewRequest(http.MethodGet, "/order/id-1", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

