package orders

import (
	"context"
	"time"
)

var Statuses = []string{"ordered", "paid", "delivered", "inhouse"}

// Metadata to be used for customized
// describing of particular Order.
type Metadata map[string]interface{}

// Order this represents the order to be made
type Order struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Price     uint64    `json:"price,omitempty"`
	Metadata  Metadata  `json:"metadata,omitempty"`
	Status    string    `json:"status,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// OrderService.
type OrderService interface {
	// Save
	CreateOrder(ctx context.Context, token string, order Order) (string, error)

	// ViewOrder retrieves Order by its unique identifier ID.
	ViewOrder(ctx context.Context, token string, id string) (Order, error)

	// ListOrders retrieves all orders for a give pageMetadata.
	ListOrders(ctx context.Context, token string, pm PageMetadata) (OrdersPage, error)

	// UpdateOrder updates the name, prices, metadata and status
	//for a given order by its unique identifier ID.
	UpdateOrder(ctx context.Context, token string, p Order) (string, error)

	// DeleteOrder deletes the order
	DeleteOrder(ctx context.Context, token string, id string) error
}

// OrderRepository specifies an account persistence API.
type OrderRepository interface {
	// Save persists the Order. A non-nil error is returned to indicate
	// operation failure.
	Save(ctx context.Context, order Order) (string, error)

	// RetrieveByID retrieves Order by its unique identifier ID.
	RetrieveByID(ctx context.Context, id string) (Order, error)

	// RetrieveAll retrieves all orders for a give pageMetadata.
	RetrieveAll(ctx context.Context, pm PageMetadata) (OrdersPage, error)

	// Update updates the name, prices, metadata and status
	//for a given order by its unique identifier ID.
	Update(ctx context.Context, p Order) (string, error)

	// Delete deletes the order
	Delete(ctx context.Context, id string) error
}

// Validate returns an error if order representation is invalid.
func (order Order) Validate() error {
	if ok := ValidateStatus(order.Status); !ok {
		return ErrInvalidStatus
	}
	return nil
}

// ValidateStatus check if the order status is acceptable
func ValidateStatus(order string) bool {
	for _, status := range Statuses {
		if status == order {
			return true
		}
	}
	return false
}
