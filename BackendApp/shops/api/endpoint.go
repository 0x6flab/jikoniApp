package api

import (
	"context"

	"github.com/0x6flab/jikoniApp/BackendApp/shops"
	"github.com/go-kit/kit/endpoint"
)

func createShopEndpoint(svc shops.ShopService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createShopReq)
		if err := req.validate(); err != nil {
			return createShopRes{}, err
		}
		sid, err := svc.CreateShop(ctx, req.token, req.shop)
		if err != nil {
			return createShopRes{}, err
		}
		ucr := createShopRes{
			ID:      sid,
			created: true,
		}

		return ucr, nil
	}
}

func viewShopEndpoint(svc shops.ShopService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewShopReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		shop, err := svc.ViewShop(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}
		return viewShopRes{
			ID:        shop.ID,
			Name:      shop.Name,
			Email:     shop.Email,
			Number:    shop.Number,
			Metadata:  shop.Metadata,
			CreatedAt: shop.CreatedAt,
			UpdatedAt: shop.UpdatedAt,
		}, nil
	}
}

func listShopsEndpoint(svc shops.ShopService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listShopsReq)
		if err := req.validate(); err != nil {
			return shops.ShopsPage{}, err
		}
		pm := shops.PageMetadata{
			Offset: req.offset,
			Limit:  req.limit,
			Total:  req.total,
			Name:   req.name,
			Email:  req.email,
			Number: req.number,
		}
		up, err := svc.ListShops(ctx, req.token, pm)
		if err != nil {
			return shops.ShopsPage{}, err
		}
		return buildShopsResponse(up), nil
	}
}

func updateShopEndpoint(svc shops.ShopService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateShopReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		shop := shops.Shop{
			ID:       req.id,
			Name:     req.Name,
			Email:    req.Email,
			Number:   req.Number,
			Metadata: req.Metadata,
		}
		oid, err := svc.UpdateShop(ctx, req.token, shop)
		if err != nil {
			return nil, err
		}
		return updateShopRes{ID: oid, updated: true}, nil
	}
}

func deleteShopEndpoint(svc shops.ShopService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteShopReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		if err := svc.DeleteShop(ctx, req.token, req.id); err != nil {
			return nil, err
		}
		return deleteShopRes{ID: req.id, deleted: true}, nil
	}
}

func buildShopsResponse(op shops.ShopsPage) shopsPageRes {
	res := shopsPageRes{
		pageRes: pageRes{
			Total:  op.Total,
			Offset: op.Offset,
			Limit:  op.Limit,
		},
		Shops: []viewShopRes{},
	}
	for _, shop := range op.Shops {
		view := viewShopRes{
			ID:        shop.ID,
			Name:      shop.Name,
			Email:     shop.Email,
			Number:    shop.Number,
			Metadata:  shop.Metadata,
			CreatedAt: shop.CreatedAt,
			UpdatedAt: shop.UpdatedAt,
		}
		res.Shops = append(res.Shops, view)
	}
	return res
}
