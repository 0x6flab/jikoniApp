package shops

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
	Name     string
	Email    string
	Number   string
	Metadata Metadata
}

// ShopsPage contains a page of shops.
type ShopsPage struct {
	PageMetadata
	Shops []Shop
}

var _ ShopService = (*shopService)(nil)

type shopService struct {
	shops ShopRepository
}

// NewShopService instantiates the users service implementation
func NewShopService(shops ShopRepository) ShopService {
	return &shopService{
		shops: shops,
	}
}

func (svc shopService) CreateShop(ctx context.Context, token string, shop Shop) (string, error) {
	if err := shop.Validate(); err != nil {
		return "", err
	}
	shop.ID = ulid.Make().String()
	shop.CreatedAt = time.Now()
	shop.UpdatedAt = time.Now()
	uid, err := svc.shops.Save(ctx, shop)
	if err != nil {
		return "", err
	}
	return uid, nil
}

func (svc shopService) ViewShop(ctx context.Context, token, id string) (Shop, error) {
	return svc.shops.RetrieveByID(ctx, id)
}

func (svc shopService) ListShops(ctx context.Context, token string, pm PageMetadata) (ShopsPage, error) {
	return svc.shops.RetrieveAll(ctx, pm)
}

func (svc shopService) UpdateShop(ctx context.Context, token string, shop Shop) (string, error) {
	uShop := Shop{
		ID:        shop.ID,
		Name:      shop.Name,
		Email:     shop.Email,
		Number:    shop.Number,
		Metadata:  shop.Metadata,
		UpdatedAt: time.Now(),
	}
	return svc.shops.Update(ctx, uShop)
}

func (svc shopService) DeleteShop(ctx context.Context, token string, id string) error {
	return svc.shops.Delete(ctx, id)
}
