package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/0x6flab/jikoniApp/BackendApp/shops"
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
	_ Response = (*createShopRes)(nil)
	_ Response = (*viewShopRes)(nil)
	_ Response = (*shopsPageRes)(nil)
	_ Response = (*updateShopRes)(nil)
	_ Response = (*deleteShopRes)(nil)
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

type createShopRes struct {
	ID      string
	created bool
}

func (res createShopRes) Code() int {
	if res.created {
		return http.StatusCreated
	}
	return http.StatusOK
}

func (res createShopRes) Headers() map[string]string {
	if res.created {
		return map[string]string{
			"Location": fmt.Sprintf("/shops/%s", res.ID),
		}
	}
	return map[string]string{}
}

func (res createShopRes) Empty() bool {
	return true
}

type viewShopRes struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email,omitempty"`
	Number    string         `json:"number,omitempty"`
	Metadata  shops.Metadata `json:"metadata,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
}

func (res viewShopRes) Code() int {
	return http.StatusOK
}

func (res viewShopRes) Headers() map[string]string {
	return map[string]string{}
}

func (res viewShopRes) Empty() bool {
	return false
}

type shopsPageRes struct {
	pageRes
	Shops []viewShopRes `json:"shops"`
}

func (res shopsPageRes) Code() int {
	return http.StatusOK
}

func (res shopsPageRes) Headers() map[string]string {
	return map[string]string{}
}

func (res shopsPageRes) Empty() bool {
	return false
}

type updateShopRes struct {
	ID      string
	updated bool
}

func (res updateShopRes) Code() int {
	return http.StatusOK
}

func (res updateShopRes) Headers() map[string]string {
	if res.updated {
		return map[string]string{
			"Location": fmt.Sprintf("/shops/%s", res.ID),
		}
	}
	return map[string]string{}
}

func (res updateShopRes) Empty() bool {
	return true
}

type deleteShopRes struct {
	ID      string
	deleted bool
}

func (res deleteShopRes) Code() int {
	return http.StatusNoContent
}

func (res deleteShopRes) Headers() map[string]string {
	if res.deleted {
		return map[string]string{
			"Location": fmt.Sprintf("/shops/%s", res.ID),
		}
	}
	return map[string]string{}
}

func (res deleteShopRes) Empty() bool {
	return true
}
