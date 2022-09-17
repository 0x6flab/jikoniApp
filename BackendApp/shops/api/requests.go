package api

import (
	"github.com/0x6flab/jikoniApp/BackendApp/internal/errors"
	"github.com/0x6flab/jikoniApp/BackendApp/shops"
)

const (
	maxLimitSize = 100
)

type createShopReq struct {
	shop  shops.Shop
	token string
}

func (req createShopReq) validate() error {
	if req.token == "" {
		return errors.ErrBearerToken
	}
	return req.shop.Validate()
}

type viewShopReq struct {
	token string
	id    string
}

func (req viewShopReq) validate() error {
	if req.token == "" {
		return errors.ErrBearerToken
	}
	if req.id == "" {
		return errors.ErrMissingID
	}
	return nil
}

type listShopsReq struct {
	token  string
	name   string
	email  string
	number string
	offset uint64
	limit  uint64
	total  uint64
}

func (req listShopsReq) validate() error {
	if req.token == "" {
		return errors.ErrBearerToken
	}
	if req.limit > maxLimitSize || req.limit < 1 {
		return errors.ErrLimitSize
	}
	return nil
}

type updateShopReq struct {
	token    string
	id       string
	Name     string         `json:"name,omitempty"`
	Email    string         `json:"email,omitempty"`
	Number   string         `json:"number,omitempty"`
	Metadata shops.Metadata `json:"metadata,omitempty"`
}

func (req updateShopReq) validate() error {
	if req.token == "" {
		return errors.ErrBearerToken
	}
	if req.id == "" {
		return errors.ErrMissingID
	}
	return nil
}

type deleteShopReq struct {
	token string
	id    string
}

func (req deleteShopReq) validate() error {
	if req.token == "" {
		return errors.ErrBearerToken
	}
	if req.id == "" {
		return errors.ErrMissingID
	}
	return nil
}
