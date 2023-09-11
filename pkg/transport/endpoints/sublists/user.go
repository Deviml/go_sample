package sublists

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/auth/jwt/claims"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	http2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/errors/http"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
)

type Endpoints struct {
	Create endpoint.Endpoint
	Show   endpoint.Endpoint
	Update endpoint.Endpoint
	Delete endpoint.Endpoint
}

type SublistService interface {
	CreateSublistForUser(ctx context.Context, userID string, createRequest CreateRequest) error
	UpdateSublist(ctx context.Context, updateRequest UpdateRequest) error
	ShowSublist(ctx context.Context, sublistID string) (entities.CompleteSublist, error)
	DeleteSublist(ctx context.Context, sublistID string) error
}

func NewEndpoints(svc SublistService) *Endpoints {
	return &Endpoints{
		Create: MakeCreateEndpoint(svc),
		Show:   MakeShowEndpoint(svc),
		Update: MakeUpdateEndpoint(svc),
		Delete: MakeDeleteEndpoint(svc),
	}
}

func MakeDeleteEndpoint(svc SublistService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		dr := request.(*SublistRequest)
		err = svc.DeleteSublist(ctx, dr.ID)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeShowEndpoint(svc SublistService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cr := request.(*SublistRequest)
		s, err := svc.ShowSublist(ctx, cr.ID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{
			Data: s,
		}, nil
	}
}

func MakeUpdateEndpoint(svc SublistService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cr := request.(*UpdateRequest)
		err = svc.UpdateSublist(ctx, *cr)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeCreateEndpoint(svc SublistService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cr := request.(*CreateRequest)
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		err = svc.CreateSublistForUser(ctx, userID, *cr)
		if err != nil {
			return nil, err
		}
		return responses.Created{}, nil
	}
}

type CompanyRequest struct {
	Category string `json:"category"`
	Name     string `json:"name"`
	ID       string `json:"id"`
}

type CreateRequest struct {
	LocationID string           `json:"-"`
	Companies  []CompanyRequest `json:"companies"`
	Name       string           `json:"name"`
	Zipcode    string           `json:"zipcode"`
	CountyID   string           `json:"county_id"`
	CityID     string           `json:"city_id"`
}

type SublistRequest struct {
	ID string `json:"id"`
}

type UpdateRequest struct {
	ID        string           `json:"-"`
	Companies []CompanyRequest `json:"companies"`
	Zipcode   string           `json:"zipcode"`
	Name      string           `json:"name"`
}

func DecodeShowRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	return &SublistRequest{ID: id}, nil
}

func DecodeCreateRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var createRequest CreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &createRequest, nil
}

func DecodePatchRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	var updateRequest UpdateRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	updateRequest.ID = id
	return &updateRequest, nil
}

func DecodeDeleteRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	return &SublistRequest{ID: id}, nil
}

type GetUserSubService interface {
	Get(ctx context.Context, userID string) ([]entities.CompleteSublist, error)
}

func MakeUserSubEndpoint(svc GetUserSubService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		s, err := svc.Get(ctx, userID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{Data: s}, nil
	}
}
