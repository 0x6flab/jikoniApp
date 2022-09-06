package orders

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total    uint64
	Offset   uint64
	Limit    uint64
	Name     string
	Price    uint64
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
	uid, err := svc.generateUUID()
	if err != nil {
		return "", err
	}
	order.ID = uid
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	uid, err = svc.orders.Save(ctx, order)
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
		Metadata:  order.Metadata,
		UpdatedAt: time.Now(),
	}
	return svc.orders.Update(ctx, uOrder)
}

func (svc orderService) DeleteOrder(ctx context.Context, token string, id string) error {
	return svc.orders.Delete(ctx, id)
}

func (svc orderService) generateUUID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
