package orders

import (
	"context"
	"time"

	"github.com/0x6flab/jikoniApp/BackendApp/internal/errors"
)

// Places describes where the order was placed or is being taken.
var Places = []string{"inhouse", "delivery"}

// Statuses describe the payment status of the order.
var Statuses = []string{"ordered", "paid"}

// Metadata to be used for customized
// describing of particular Order.
type Metadata map[string]interface{}

// Order this represents the order to be made by a person to the shop.
type Order struct {
	ID        string    `json:"id,omitempty"`
	Vendor    string    `json:"vendor,omitempty"`     // The ID of the vendor of the product i.e shop.
	Name      string    `json:"name,omitempty"`       // The name of the order good.
	Price     uint64    `json:"price,omitempty"`      // This is the price of the order.
	Place     string    `json:"place,omitempty"`      // This is the place where the order was served. It is either inhouse or delivery.
	Status    string    `json:"status,omitempty"`     // This is the payment status. It is either paid or ordered.
	Metadata  Metadata  `json:"metadata,omitempty"`   // Metadata contains extra information about the order.
	UpdatedAt time.Time `json:"updated_at,omitempty"` // When the order was updated.
	CreatedAt time.Time `json:"created_at,omitempty"` // When the order was created in the system.
}

// OrderService. This describes the methods an Order undergo.
// CreateOrder
// ViewOrder
// ListOrders
// UpdateOrder
// DeleteOrder
type OrderService interface {
	// CreateOrder creates and order to the system. Requires a token and the order object.
	CreateOrder(ctx context.Context, token string, order Order) (string, error)

	// ViewOrder retrieves Order by its unique identifier ID.
	ViewOrder(ctx context.Context, token string, id string) (Order, error)

	// ListOrders retrieves all orders for a give pageMetadata.
	ListOrders(ctx context.Context, token string, pm PageMetadata) (OrdersPage, error)

	// UpdateOrder updates the name, prices, metadata, place and status
	// for a given order by its unique identifier ID.
	UpdateOrder(ctx context.Context, token string, p Order) (string, error)

	// DeleteOrder deletes the order for a give unique identifier ID.
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

	// Update updates the name, prices, metadata, place and status
	// for a given order by its unique identifier ID.
	Update(ctx context.Context, p Order) (string, error)

	// Delete deletes the order
	Delete(ctx context.Context, id string) error
}

// Validate returns an error if order representation is invalid.
func (order Order) Validate() error {
	if ok := ValidateStatus(order.Status); !ok {
		return errors.ErrInvalidStatus
	}
	if ok := ValidatePlaces(order.Place); !ok {
		return errors.ErrInvalidStatus
	}
	return nil
}

// ValidatePlaces check if the order place is acceptable
func ValidatePlaces(order string) bool {
	for _, place := range Places {
		if place == order {
			return true
		}
	}
	return false
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
