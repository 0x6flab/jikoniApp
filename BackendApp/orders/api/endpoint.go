package api

import (
	"context"

	"github.com/0x6flab/jikoniApp/BackendApp/orders"
	"github.com/go-kit/kit/endpoint"
)

func createOrderEndpoint(svc orders.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createOrderReq)
		if err := req.validate(); err != nil {
			return createOrderRes{}, err
		}
		uid, err := svc.CreateOrder(ctx, req.token, req.order)
		if err != nil {
			return createOrderRes{}, err
		}
		ucr := createOrderRes{
			ID:      uid,
			created: true,
		}

		return ucr, nil
	}
}

func viewOrderEndpoint(svc orders.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(viewOrderReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		order, err := svc.ViewOrder(ctx, req.token, req.id)
		if err != nil {
			return nil, err
		}
		return viewOrderRes{
			ID:        order.ID,
			Name:      order.Name,
			Price:     order.Price,
			Metadata:  order.Metadata,
			Status:    order.Status,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		}, nil
	}
}

func listOrdersEndpoint(svc orders.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listOrdersReq)
		if err := req.validate(); err != nil {
			return orders.OrdersPage{}, err
		}
		pm := orders.PageMetadata{
			Offset: req.offset,
			Limit:  req.limit,
			Total:  req.total,
			Name:   req.name,
			Price:  req.price,
			Status: req.status,
		}
		up, err := svc.ListOrders(ctx, req.token, pm)
		if err != nil {
			return orders.OrdersPage{}, err
		}
		return buildOrdersResponse(up), nil
	}
}

func updateOrderEndpoint(svc orders.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateOrderReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		order := orders.Order{
			ID:       req.id,
			Name:     req.Name,
			Price:    req.Price,
			Status:   req.Status,
			Metadata: req.Metadata,
		}
		oid, err := svc.UpdateOrder(ctx, req.token, order)
		if err != nil {
			return nil, err
		}
		return updateOrderRes{ID: oid, updated: true}, nil
	}
}

func deleteOrderEndpoint(svc orders.OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteOrderReq)
		if err := req.validate(); err != nil {
			return nil, err
		}
		if err := svc.DeleteOrder(ctx, req.token, req.id); err != nil {
			return nil, err
		}
		return deleteOrderRes{ID: req.id, deleted: true}, nil
	}
}

func buildOrdersResponse(op orders.OrdersPage) ordersPageRes {
	res := ordersPageRes{
		pageRes: pageRes{
			Total:  op.Total,
			Offset: op.Offset,
			Limit:  op.Limit,
		},
		Orders: []viewOrderRes{},
	}
	for _, order := range op.Orders {
		view := viewOrderRes{
			ID:        order.ID,
			Name:      order.Name,
			Price:     order.Price,
			Status:    order.Status,
			Metadata:  order.Metadata,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		}
		res.Orders = append(res.Orders, view)
	}
	return res
}
