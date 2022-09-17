//go:build !test

package api

import (
	"context"
	"time"

	"github.com/0x6flab/jikoniApp/BackendApp/shops"
	"github.com/go-kit/kit/metrics"
)

var _ shops.ShopService = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     shops.ShopService
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc shops.ShopService, counter metrics.Counter, latency metrics.Histogram) shops.ShopService {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (ms *metricsMiddleware) CreateShop(ctx context.Context, token string, shop shops.Shop) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "create_shop").Add(1)
		ms.latency.With("method", "create_shop").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.CreateShop(ctx, token, shop)
}
func (ms *metricsMiddleware) ViewShop(ctx context.Context, token, id string) (shops.Shop, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "view_shop").Add(1)
		ms.latency.With("method", "view_shop").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ViewShop(ctx, token, id)
}
func (ms *metricsMiddleware) ListShops(ctx context.Context, token string, pm shops.PageMetadata) (shops.ShopsPage, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "list_shops").Add(1)
		ms.latency.With("method", "list_shops").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.ListShops(ctx, token, pm)
}
func (ms *metricsMiddleware) UpdateShop(ctx context.Context, token string, shop shops.Shop) (string, error) {
	defer func(begin time.Time) {
		ms.counter.With("method", "update_shop").Add(1)
		ms.latency.With("method", "update_shop").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.UpdateShop(ctx, token, shop)
}
func (ms *metricsMiddleware) DeleteShop(ctx context.Context, token, id string) error {
	defer func(begin time.Time) {
		ms.counter.With("method", "delete_shop").Add(1)
		ms.latency.With("method", "delete_shop").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return ms.svc.DeleteShop(ctx, token, id)
}
