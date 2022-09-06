package api

import (
	"github.com/0x6flab/jikoniApp/BackendApp/internal/errors"
	"github.com/0x6flab/jikoniApp/BackendApp/orders"
)

const (
	maxLimitSize = 100
)

type createOrderReq struct {
	order orders.Order
	token string
}

func (req createOrderReq) validate() error {
	if req.token == "" {
		return errors.ErrBearerToken
	}
	return req.order.Validate()
}

type viewOrderReq struct {
	token string
	id    string
}

func (req viewOrderReq) validate() error {
	if req.token == "" {
		return errors.ErrBearerToken
	}
	if req.id == "" {
		return errors.ErrMissingID
	}
	return nil
}

type listOrdersReq struct {
	token  string
	name   string
	price  uint64
	status string
	offset uint64
	limit  uint64
	total  uint64
}

func (req listOrdersReq) validate() error {
	if req.token == "" {
		return errors.ErrBearerToken
	}
	if req.limit > maxLimitSize || req.limit < 1 {
		return errors.ErrLimitSize
	}
	return nil
}

type updateOrderReq struct {
	token    string
	id       string
	Name     string          `json:"name,omitempty"`
	Price    uint64          `json:"price,omitempty"`
	Status   string          `json:"status,omitempty"`
	Metadata orders.Metadata `json:"metadata,omitempty"`
}

func (req updateOrderReq) validate() error {
	if req.token == "" {
		return errors.ErrBearerToken
	}
	if req.id == "" {
		return errors.ErrMissingID
	}
	return nil
}

type deleteOrderReq struct {
	token string
	id    string
}

func (req deleteOrderReq) validate() error {
	if req.token == "" {
		return errors.ErrBearerToken
	}
	if req.id == "" {
		return errors.ErrMissingID
	}
	return nil
}
