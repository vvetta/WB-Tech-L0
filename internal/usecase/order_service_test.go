package usecase_test

import (
	"errors"
	"context"
	"testing"

	"WB-Tech-L0/internal/domain"
	"WB-Tech-L0/internal/usecase"
	"WB-Tech-L0/internal/usecase/mocks"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type noopLogger struct{}
func (noopLogger) Info(string, ...any) {}
func (noopLogger) Error(string, ...any) {}
func (noopLogger) Debug(string, ...any) {}

func fakeOrder() *domain.Order {
	gofakeit.Seed(0)
	return &domain.Order{
		OrderUID: gofakeit.UUID(),
		TrackNumber: gofakeit.LetterN(10),
		Entry: "WBIL",
		Delivery: domain.DeliveryInfo{
			Name: gofakeit.Name(),
			Phone: gofakeit.Phone(),
			Zip: gofakeit.Zip(),
			City: gofakeit.City(),
			Address: gofakeit.Street(),
			Region: gofakeit.State(),
			Email: gofakeit.Email(),
		},
		Payment: domain.PaymentInfo{
			Transaction: gofakeit.UUID(),
			Currency: "USD",
			Provider: "wbpay",
			Amount: 100,
			PaymentDT: 170000,
			Bank: "alpha",
			DeliveryCost: 10,
			GoodsTotal: 90,
		},
		Locale: "en",
		CustomerID: "test",
		DeliveryService: "meest",
		ShardKey: "9",
		SMID: 99,
		DateCreated: "2025",
		OOFShard: "1",
	}
}

func TestOrderService_GetByID_emptyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepo(ctrl)
	cache := mocks.NewMockCache(ctrl)

	svc := usecase.NewOrderService(repo, cache, noopLogger{})
	_, err := svc.GetByID(context.Background(), "")
	require.Error(t, err)
}

func TestOrderService_GetByID_cacheHit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepo(ctrl)
	cache := mocks.NewMockCache(ctrl)

	o := fakeOrder()
	cache.EXPECT().Get(o.OrderUID).Return(o, true)

	svc := usecase.NewOrderService(repo, cache, noopLogger{})
	got, err := svc.GetByID(context.Background(), o.OrderUID)
	require.NoError(t, err)
	require.Equal(t, o, got)
}

func TestOrderService_GetByID_cacheMiss(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepo(ctrl)
	cache := mocks.NewMockCache(ctrl)

	o := fakeOrder()
	cache.EXPECT().Get(o.OrderUID).Return(nil, false)
	repo.EXPECT().GetOrderById(gomock.Any(), o.OrderUID).Return(o, nil)
	cache.EXPECT().Set(o.OrderUID, o)

	svc := usecase.NewOrderService(repo, cache, noopLogger{})
	got, err := svc.GetByID(context.Background(), o.OrderUID)
	require.NoError(t, err)
	require.Equal(t, o, got)
}

func TestOrderService_GetByID_notFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepo(ctrl)
	cache := mocks.NewMockCache(ctrl)

	id := "missing"
	cache.EXPECT().Get(id).Return(nil, false)
	repo.EXPECT().GetOrderById(gomock.Any(), id).Return(nil, domain.ErrNotFound)

	svc := usecase.NewOrderService(repo, cache, noopLogger{})
	got, err := svc.GetByID(context.Background(), id)
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrNotFound)
	require.Nil(t, got)
}

func TestOrderService_WarmUpCache_ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepo(ctrl)
	cache := mocks.NewMockCache(ctrl)

	o1 := fakeOrder()
	o2 := fakeOrder()

	repo.EXPECT().ListRecentOrders(gomock.Any(), 5).Return([]*domain.Order{o1, o2}, nil)
	cache.EXPECT().Set(o1.OrderUID, o1)
	cache.EXPECT().Set(o2.OrderUID, o2)

	svc := usecase.NewOrderService(repo, cache, noopLogger{})
	err := svc.WarmUpCache(context.Background(), 5)
	require.NoError(t, err)
}

func TestOrderService_WarmUpCache_empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepo(ctrl)
	cache := mocks.NewMockCache(ctrl)

	repo.EXPECT().ListRecentOrders(gomock.Any(), 3).Return([]*domain.Order{}, nil)

	svc := usecase.NewOrderService(repo, cache, noopLogger{})
	err := svc.WarmUpCache(context.Background(), 3)
	require.NoError(t, err)
}

func TestOrderService_WarmUpCache_dbError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockRepo(ctrl)
	cache := mocks.NewMockCache(ctrl)

	dbErr := errors.New("db down")
	repo.EXPECT().ListRecentOrders(gomock.Any(), 3).Return(nil, dbErr)
	
	svc := usecase.NewOrderService(repo, cache, noopLogger{})
	err := svc.WarmUpCache(context.Background(), 3)
	require.Error(t, err)
	require.ErrorIs(t, err, dbErr)
}

