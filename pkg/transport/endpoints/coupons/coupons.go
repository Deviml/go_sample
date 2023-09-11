package coupons2

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/auth/jwt/claims"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	http2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/errors/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	ValidateCoupon endpoint.Endpoint
	UpdateCoupon   endpoint.Endpoint
}

type CouponService interface {
	ValidateCoupon(request ValidateCouponRequest, userID string) (entities.Coupon, error)
	UpdateCoupon(request UpdateCouponRequest, userID string) (string, error)
}

func NewEndpoints(s CouponService) *Endpoints {
	return &Endpoints{
		ValidateCoupon: makeValidateCouponEndpoint(s),
		UpdateCoupon:   makeUpdateCouponEndpoint(s),
	}
}

func makeValidateCouponEndpoint(s CouponService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*ValidateCouponRequest)
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token == nil {
			return nil, http2.BadRequest{}
		}

		userToken, ok := token.(*claims.UserClaims)
		if !ok {
			return nil, http2.BadRequest{}
		}
		valid, err := s.ValidateCoupon(*req, userToken.UserID)
		if err != nil {
			return nil, err
		}
		return valid, nil
	}
}

type ValidateCouponRequest struct {
	CouponName string `json:"coupon_name"`
}

func DecodeValidateCouponRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var p ValidateCouponRequest
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &p, nil
}

func makeUpdateCouponEndpoint(s CouponService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*UpdateCouponRequest)
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token == nil {
			return nil, http2.BadRequest{}
		}

		userToken, ok := token.(*claims.UserClaims)
		if !ok {
			return nil, http2.BadRequest{}
		}
		valid, err := s.UpdateCoupon(*req, userToken.UserID)
		if err != nil {
			return nil, err
		}
		return valid, nil
	}
}

type UpdateCouponRequest struct {
	CouponName string `json:"coupon_name"`
}

func DecodeUpdateCouponRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var p UpdateCouponRequest
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &p, nil
}
