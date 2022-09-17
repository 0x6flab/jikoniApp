//go:build !test

package api

import (
	"context"
	"time"

	"github.com/0x6flab/jikoniApp/BackendApp/shops"
	"github.com/go-kit/log"
)

var _ shops.ShopService = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    shops.ShopService
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc shops.ShopService, logger log.Logger) shops.ShopService {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) CreateShop(ctx context.Context, token string, shop shops.Shop) (id string, err error) {
	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "create_shop",
			"token", token,
			"name", shop.Name,
			"number", shop.Number,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return lm.svc.CreateShop(ctx, token, shop)

}
func (lm *loggingMiddleware) ViewShop(ctx context.Context, token, id string) (shop shops.Shop, err error) {

	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "view_shop",
			"token", token,
			"id", id,
			"name", shop.Name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return lm.svc.ViewShop(ctx, token, id)

}
func (lm *loggingMiddleware) ListShops(ctx context.Context, token string, pm shops.PageMetadata) (fp shops.ShopsPage, err error) {
	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "list_shops",
			"token", token,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return lm.svc.ListShops(ctx, token, pm)

}
func (lm *loggingMiddleware) UpdateShop(ctx context.Context, token string, shop shops.Shop) (rid string, err error) {
	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "update_shop",
			"token", token,
			"id", shop.ID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return lm.svc.UpdateShop(ctx, token, shop)

}
func (lm *loggingMiddleware) DeleteShop(ctx context.Context, token, id string) (err error) {

	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "delete_shop",
			"token", token,
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return lm.svc.DeleteShop(ctx, token, id)

}
