package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/0x6flab/jikoniApp/BackendApp/orders"
)

// Response contains HTTP response specific methods.
type Response interface {
	// Code returns HTTP response code.
	Code() int

	// Headers returns map of HTTP headers with their values.
	Headers() map[string]string

	// Empty indicates if HTTP response has content.
	Empty() bool
}

var (
	_ Response = (*tokenRes)(nil)
	_ Response = (*createOrderRes)(nil)
	_ Response = (*viewOrderRes)(nil)
	_ Response = (*ordersPageRes)(nil)
	_ Response = (*updateOrderRes)(nil)
	_ Response = (*deleteOrderRes)(nil)
)

type pageRes struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

type tokenRes struct {
	Token string `json:"token,omitempty"`
}

func (res tokenRes) Code() int {
	return http.StatusCreated
}

func (res tokenRes) Headers() map[string]string {
	return map[string]string{}
}

func (res tokenRes) Empty() bool {
	return res.Token == ""
}

type createOrderRes struct {
	ID      string
	created bool
}

func (res createOrderRes) Code() int {
	if res.created {
		return http.StatusCreated
	}
	return http.StatusOK
}

func (res createOrderRes) Headers() map[string]string {
	if res.created {
		return map[string]string{
			"Location": fmt.Sprintf("/orders/%s", res.ID),
		}
	}
	return map[string]string{}
}

func (res createOrderRes) Empty() bool {
	return true
}

type viewOrderRes struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Price     uint64          `json:"price,omitempty"`
	Metadata  orders.Metadata `json:"metadata,omitempty"`
	Status    string          `json:"status,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
}

func (res viewOrderRes) Code() int {
	return http.StatusOK
}

func (res viewOrderRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewOrderRes) Empty() bool {
	return false
}

type ordersPageRes struct {
	pageRes
	Orders []viewOrderRes `json:"orders"`
}

func (res ordersPageRes) Code() int {
	return http.StatusOK
}

func (res ordersPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res ordersPageRes) Empty() bool {
	return false
}

type updateOrderRes struct {
	ID      string
	updated bool
}

func (res updateOrderRes) Code() int {
	return http.StatusOK
}

func (res updateOrderRes) Headers() map[string]string {
	if res.updated {
		return map[string]string{
			"Location": fmt.Sprintf("/orders/%s", res.ID),
		}
	}
	return map[string]string{}
}

func (res updateOrderRes) Empty() bool {
	return true
}

type deleteOrderRes struct {
	ID      string
	deleted bool
}

func (res deleteOrderRes) Code() int {
	return http.StatusNoContent
}

func (res deleteOrderRes) Headers() map[string]string {
	if res.deleted {
		return map[string]string{
			"Location": fmt.Sprintf("/orders/%s", res.ID),
		}
	}
	return map[string]string{}
}

func (res deleteOrderRes) Empty() bool {
	return true
}
