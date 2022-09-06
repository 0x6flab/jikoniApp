//go:build !test

package api

import (
	"context"
	"time"

	"github.com/0x6flab/jikoniApp/BackendApp/orders"
	"github.com/go-kit/kit/metrics"
)

var _ orders.OrderService = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     orders.OrderService
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc orders.OrderService, counter metrics.Counter, latency metrics.Histogram) orders.OrderService {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) CreateOrder(ctx context.Context, token string, order orders.Order) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_order").Add(1)
		ms.latency.With("method", "create_order").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateOrder(ctx, token, order)
}
func (ms *metricsMiddleware) ViewOrder(ctx context.Context, token, id string) (orders.Order, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_order").Add(1)
		ms.latency.With("method", "view_order").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewOrder(ctx, token, id)
}
func (ms *metricsMiddleware) ListOrders(ctx context.Context, token string, pm orders.PageMetadata) (orders.OrdersPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_orders").Add(1)
		ms.latency.With("method", "list_orders").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListOrders(ctx, token, pm)
}
func (ms *metricsMiddleware) UpdateOrder(ctx context.Context, token string, order orders.Order) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_order").Add(1)
		ms.latency.With("method", "update_order").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateOrder(ctx, token, order)
}
func (ms *metricsMiddleware) DeleteOrder(ctx context.Context, token, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "delete_order").Add(1)
		ms.latency.With("method", "delete_order").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.DeleteOrder(ctx, token, id)
}
