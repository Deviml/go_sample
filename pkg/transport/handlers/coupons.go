package handlers

import (
	"net/http"

	coupons2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/coupons"
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const (
	ValidateCouponPrefix = "/validate-coupon"
	ValidateCoupon       = ""
	UpdateCouponPrefix   = "/update-coupon"
	UpdateCoupon         = ""
)

func MakeCouponsHandler(r *mux.Router, eps *coupons2.Endpoints, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	addValidateCouponRoute(r, authMiddleware(eps.ValidateCoupon), op)
	addUpdateCouponRoute(r, authMiddleware(eps.UpdateCoupon), op)
	return r
}

func addValidateCouponRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(ValidateCouponPrefix).
		Path(ValidateCoupon).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			coupons2.DecodeValidateCouponRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addUpdateCouponRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(UpdateCouponPrefix).
		Path(UpdateCoupon).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			coupons2.DecodeUpdateCouponRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}
