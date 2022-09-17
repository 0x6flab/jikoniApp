package shops

import (
	"context"
	"time"

	"github.com/0x6flab/jikoniApp/BackendApp/internal/errors"
)

// Metadata to be used for customized
// describing of particular Shop.
type Metadata map[string]interface{}

// Shop this represents the shop where the order is being made.
type Shop struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`       // The name of the shop.
	Email     string    `json:"email,omitempty"`      // This is the email of the shop.
	Number    string    `json:"number,omitempty"`     // This is the phone number of the shop.
	Metadata  Metadata  `json:"metadata,omitempty"`   // Metadata contains extra information about the shop.
	UpdatedAt time.Time `json:"updated_at,omitempty"` // When the shop was updated.
	CreatedAt time.Time `json:"created_at,omitempty"` // When the shop was created in the system.
}

// ShopService. This describes the methods a Shop undergo.
// CreateShop
// ViewShop
// ListShops
// UpdateShop
// DeleteShop
type ShopService interface {
	// CreateShop creates a shop to the system. Requires a token and the shop object.
	CreateShop(ctx context.Context, token string, shop Shop) (string, error)

	// ViewShop retrieves Shop by its unique identifier ID.
	ViewShop(ctx context.Context, token string, id string) (Shop, error)

	// ListShops retrieves all shops for a give pageMetadata.
	ListShops(ctx context.Context, token string, pm PageMetadata) (ShopsPage, error)

	// UpdateShop updates the name, email, phone number and metadata
	// for a given shop by its unique identifier ID.
	UpdateShop(ctx context.Context, token string, shop Shop) (string, error)

	// DeleteShop deletes the order for a give unique identifier ID.
	DeleteShop(ctx context.Context, token string, id string) error
}

// ShopRepository specifies an account persistence API.
type ShopRepository interface {
	// Save persists the Order. A non-nil error is returned to indicate
	// operation failure.
	Save(ctx context.Context, shop Shop) (string, error)

	// RetrieveByID retrieves Order by its unique identifier ID.
	RetrieveByID(ctx context.Context, id string) (Shop, error)

	// RetrieveAll retrieves all orders for a give pageMetadata.
	RetrieveAll(ctx context.Context, pm PageMetadata) (ShopsPage, error)

	// Update updates the name, email, phone number and metadata
	// for a given order by its unique identifier ID.
	Update(ctx context.Context, shop Shop) (string, error)

	// Delete deletes the shop
	Delete(ctx context.Context, id string) error
}

// Validate returns an error if order representation is invalid.
func (shop Shop) Validate() error {
	if ok := ValidateNumber(shop.Number); !ok {
		return errors.ErrInvalidStatus
	}
	if ok := ValidateEmail(shop.Email); !ok {
		return errors.ErrInvalidStatus
	}
	return nil
}

// ValidateEmail check if the shop email is acceptable
func ValidateEmail(order string) bool {
	return true
}

// ValidateNumber check if the shop phone number is acceptable
func ValidateNumber(order string) bool {
	return true
}
