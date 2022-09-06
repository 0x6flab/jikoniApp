//go:build !test

package api

import (
	"context"
	"time"

	"github.com/0x6flab/jikoniApp/BackendApp/orders"
	"github.com/go-kit/log"
)

var _ orders.OrderService = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    orders.OrderService
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc orders.OrderService, logger log.Logger) orders.OrderService {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) CreateOrder(ctx context.Context, token string, order orders.Order) (id string, err error) {
	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "create_order",
			"token", token,
			"name", order.Name,
			"price", order.Price,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return lm.svc.CreateOrder(ctx, token, order)

}
func (lm *loggingMiddleware) ViewOrder(ctx context.Context, token, id string) (order orders.Order, err error) {

	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "view_order",
			"token", token,
			"id", id,
			"name", order.Name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return lm.svc.ViewOrder(ctx, token, id)

}
func (lm *loggingMiddleware) ListOrders(ctx context.Context, token string, pm orders.PageMetadata) (fp orders.OrdersPage, err error) {
	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "list_orders",
			"token", token,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return lm.svc.ListOrders(ctx, token, pm)

}
func (lm *loggingMiddleware) UpdateOrder(ctx context.Context, token string, order orders.Order) (rid string, err error) {
	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "update_order",
			"token", token,
			"id", order.ID,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return lm.svc.UpdateOrder(ctx, token, order)

}
func (lm *loggingMiddleware) DeleteOrder(ctx context.Context, token, id string) (err error) {

	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "delete_order",
			"token", token,
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return lm.svc.DeleteOrder(ctx, token, id)

}
