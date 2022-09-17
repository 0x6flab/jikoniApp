package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/0x6flab/jikoniApp/BackendApp/internal/apiutil"
	"github.com/0x6flab/jikoniApp/BackendApp/internal/errors"
	"github.com/0x6flab/jikoniApp/BackendApp/shops"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	contentType = "application/json"
	offsetKey   = "offset"
	limitKey    = "limit"
	totalKey    = "total"
	nameKey     = "name"
	emailKey    = "email"
	numberKey   = "number"
)

// MakeShopsHandler returns a HTTP handler for API endpoints.
func MakeShopsHandler(svc shops.ShopService, r *mux.Router, logger kitlog.Logger) {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerErrorLogger(logger),
		kitoc.HTTPServerTrace(),
	}

	r.Methods("POST").Path("/shops").Handler(kithttp.NewServer(
		kitoc.TraceEndpoint("gokit:endpoint create_shop")(createShopEndpoint(svc)),
		decodeCreateShop,
		encodeResponse,
		opts...,
	))

	r.Methods("GET").Path("/shops/{id}").Handler(kithttp.NewServer(
		kitoc.TraceEndpoint("gokit:endpoint view_shop")(viewShopEndpoint(svc)),
		decodeViewShop,
		encodeResponse,
		opts...,
	))

	r.Methods("GET").Path("/shops").Handler(kithttp.NewServer(
		kitoc.TraceEndpoint("gokit:endpoint list_shops")(listShopsEndpoint(svc)),
		decodeListShops,
		encodeResponse,
		opts...,
	))

	r.Methods("PUT").Path("/shops/{id}").Handler(kithttp.NewServer(
		kitoc.TraceEndpoint("gokit:endpoint update_shop")(updateShopEndpoint(svc)),
		decodeUpdateShop,
		encodeResponse,
		opts...,
	))

	r.Methods("DELETE").Path("/shops/{id}").Handler(kithttp.NewServer(
		kitoc.TraceEndpoint("gokit:endpoint delete_shop")(deleteShopEndpoint(svc)),
		decodeDeleteShop,
		encodeResponse,
		opts...,
	))

	r.Path("/metrics").Handler(promhttp.Handler())
}

func decodeCreateShop(_ context.Context, r *http.Request) (interface{}, error) {
	if !strings.Contains(r.Header.Get("Content-Type"), contentType) {
		return nil, errors.ErrUnsupportedContentType
	}
	var shop shops.Shop
	if err := json.NewDecoder(r.Body).Decode(&shop); err != nil {
		return nil, errors.Wrap(errors.ErrMalformedEntity, err)
	}
	req := createShopReq{
		shop:  shop,
		token: decodeToken(r),
	}
	return req, nil
}

func decodeViewShop(_ context.Context, r *http.Request) (interface{}, error) {
	req := viewShopReq{
		token: decodeToken(r),
		id:    mux.Vars(r)["id"],
	}
	return req, nil
}

func decodeListShops(_ context.Context, r *http.Request) (interface{}, error) {
	var offset = uint64(0)
	var limit = uint64(100)
	var total = uint64(100)
	var name = ""
	var email = ""
	var number = ""
	var err error

	if r.URL.Query().Has(offsetKey) {
		offset, err = strconv.ParseUint(r.URL.Query().Get(offsetKey), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if r.URL.Query().Has(limitKey) {
		limit, err = strconv.ParseUint(r.URL.Query().Get(limitKey), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if r.URL.Query().Has(totalKey) {
		total, err = strconv.ParseUint(r.URL.Query().Get(totalKey), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if r.URL.Query().Has(nameKey) {
		name = r.URL.Query().Get(nameKey)
	}
	if r.URL.Query().Has(emailKey) {
		email = r.URL.Query().Get(emailKey)
	}
	if r.URL.Query().Has(numberKey) {
		number = r.URL.Query().Get(numberKey)
	}
	req := listShopsReq{
		token:  decodeToken(r),
		offset: offset,
		limit:  limit,
		total:  total,
		name:   name,
		email:  email,
		number: number,
	}
	return req, nil
}

func decodeUpdateShop(_ context.Context, r *http.Request) (interface{}, error) {
	req := updateShopReq{
		token: decodeToken(r),
		id:    mux.Vars(r)["id"],
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(errors.ErrMalformedEntity, err)
	}
	return req, nil
}

func decodeDeleteShop(_ context.Context, r *http.Request) (interface{}, error) {
	req := deleteShopReq{
		token: decodeToken(r),
		id:    mux.Vars(r)["id"],
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if ar, ok := response.(Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(ar.Code())
		if ar.Empty() {
			return nil
		}
	}
	return json.NewEncoder(w).Encode(response)
}

func decodeToken(r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	return tokenString
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch {
	case errors.Contains(err, errors.ErrInvalidQueryParams),
		errors.Contains(err, errors.ErrMalformedEntity),
		err == apiutil.ErrLimitSize,
		err == apiutil.ErrOffsetSize:
		w.WriteHeader(http.StatusBadRequest)
	case errors.Contains(err, errors.ErrAuthentication),
		err == apiutil.ErrBearerToken:
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Contains(err, errors.ErrConflict),
		errors.Contains(err, errors.ErrConflict):
		w.WriteHeader(http.StatusConflict)
	case errors.Contains(err, errors.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Contains(err, errors.ErrCreateEntity),
		errors.Contains(err, errors.ErrUpdateEntity),
		errors.Contains(err, errors.ErrViewEntity),
		errors.Contains(err, errors.ErrRemoveEntity):
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	if errorVal, ok := err.(errors.Error); ok {
		w.Header().Set("Content-Type", contentType)
		if err := json.NewEncoder(w).Encode(apiutil.ErrorRes{Err: errorVal.Msg()}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
