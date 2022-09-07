package orders

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
)

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total    uint64
	Offset   uint64
	Limit    uint64
	Vendor   string
	Name     string
	Price    uint64
	Place    string
	Metadata Metadata
	Status   string
}

// OrdersPage contains a page of orders.
type OrdersPage struct {
	PageMetadata
	Orders []Order
}

var _ OrderService = (*orderService)(nil)

type orderService struct {
	orders OrderRepository
}

// NewOrderService instantiates the users service implementation
func NewOrderService(orders OrderRepository) OrderService {
	return &orderService{
		orders: orders,
	}
}

func (svc orderService) CreateOrder(ctx context.Context, token string, order Order) (string, error) {
	if err := order.Validate(); err != nil {
		return "", err
	}
	order.ID = ulid.Make().String()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	uid, err := svc.orders.Save(ctx, order)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func (svc orderService) ViewOrder(ctx context.Context, token, id string) (Order, error) {
	return svc.orders.RetrieveByID(ctx, id)
}

func (svc orderService) ListOrders(ctx context.Context, token string, pm PageMetadata) (OrdersPage, error) {
	return svc.orders.RetrieveAll(ctx, pm)
}

func (svc orderService) UpdateOrder(ctx context.Context, token string, order Order) (string, error) {
	uOrder := Order{
		ID:        order.ID,
		Name:      order.Name,
		Price:     order.Price,
		Place:     order.Place,
		Status:    order.Status,
		Metadata:  order.Metadata,
		UpdatedAt: time.Now(),
	}
	return svc.orders.Update(ctx, uOrder)
}

func (svc orderService) DeleteOrder(ctx context.Context, token string, id string) error {
	return svc.orders.Delete(ctx, id)
}
