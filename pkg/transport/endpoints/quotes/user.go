package quotes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/auth/jwt/claims"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	http2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/errors/http"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/helpers"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
)

type Endpoints struct {
	CreateEquipment endpoint.Endpoint
	UpdateEquipment endpoint.Endpoint
	ShowEquipment   endpoint.Endpoint
	DeleteEquipment endpoint.Endpoint
	CreateSupply    endpoint.Endpoint
	UpdateSupply    endpoint.Endpoint
	ShowSupply      endpoint.Endpoint
	DeleteSupply    endpoint.Endpoint
}

type EquipmentService interface {
	CreateEquipment(ctx context.Context, userID string, createRequest CreateRequest) error
	ShowEquipment(ctx context.Context, quoteID string) (entities.CompleteEquipmentRequest, error)
	UpdateEquipment(ctx context.Context, updateRequest UpdateRequest) error
	DeleteEquipment(ctx context.Context, quoteID string) error
}

type SupplyService interface {
	CreateSupply(ctx context.Context, userID string, createRequest SupplyCreateRequest) error
	ShowSupply(ctx context.Context, quoteID string) (entities.CompleteSupplyRequest, error)
	UpdateSupply(ctx context.Context, updateRequest SupplyUpdateRequest) error
	DeleteSupply(ctx context.Context, quoteID string) error
}

func NewEndpoints(svce EquipmentService, svcs SupplyService) *Endpoints {
	return &Endpoints{
		CreateEquipment: MakeCreateEquipmentRequest(svce),
		UpdateEquipment: MakeUpdateEquipmentRequest(svce),
		ShowEquipment:   MakeShowEquipmentRequest(svce),
		DeleteEquipment: MakeDeleteEquipmentRequest(svce),
		CreateSupply:    MakeCreateSupplyRequest(svcs),
		UpdateSupply:    MakeUpdateSupplyRequest(svcs),
		ShowSupply:      MakeShowSupplyRequest(svcs),
		DeleteSupply:    MakeDeleteSupplyRequest(svcs),
	}
}

func MakeDeleteSupplyRequest(svc SupplyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		dr := request.(*ShowSupplyRequest)
		err = svc.DeleteSupply(ctx, dr.ID)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeShowSupplyRequest(svc SupplyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		sr := request.(*ShowSupplyRequest)
		s, err := svc.ShowSupply(ctx, sr.ID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{Data: s}, nil
	}
}

func MakeUpdateSupplyRequest(svc SupplyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ur := request.(*SupplyUpdateRequest)
		err = svc.UpdateSupply(ctx, *ur)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeCreateSupplyRequest(svc SupplyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cr := request.(*SupplyCreateRequest)
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		err = svc.CreateSupply(ctx, userID, *cr)
		if err != nil {
			return nil, err
		}
		return responses.Created{}, nil
	}
}

func MakeDeleteEquipmentRequest(svc EquipmentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		dr := request.(*ShowEquipmentRequest)
		err = svc.DeleteEquipment(ctx, dr.ID)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeShowEquipmentRequest(svc EquipmentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		sr := request.(*ShowEquipmentRequest)
		e, err := svc.ShowEquipment(ctx, sr.ID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{Data: e}, nil
	}
}

func MakeUpdateEquipmentRequest(svc EquipmentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ur := request.(*UpdateRequest)
		err = svc.UpdateEquipment(ctx, *ur)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeCreateEquipmentRequest(svc EquipmentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cr := request.(*CreateRequest)
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		err = svc.CreateEquipment(ctx, userID, *cr)
		if err != nil {
			return nil, err
		}
		return responses.Created{}, nil
	}
}

type CreateRequest struct {
	EquipmentID         string `json:"equipment_id"`
	EquipmentCategoryID string `json:"equipment_category_id"`
	ContractPreference  int    `json:"contract_preference"`
	StartDate           *int   `json:"start_date"`
	EndDate             *int   `json:"end_date"`
	Amount              string `json:"amount"`
	SpecialRequest      string `json:"special_request"`
	Zipcode             string `json:"zipcode"`
	Address             string `json:"address"`
	CityID              string `json:"city_id"`
	CountyID            string `json:"county_id"`
}

func DecodeCreateEquipmentRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var createRequest CreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		return nil, err
	}
	return &createRequest, nil
}

type UpdateRequest struct {
	ID                 string `json:"-"`
	ContractPreference string `json:"contract_preference"`
	Amount             int    `json:"amount"`
	SpecialRequest     string `json:"special_request"`
	Zipcode            string `json:"zipcode"`
}

func DecodeUpdateEquipmentRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	var updateRequest UpdateRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		return nil, err
	}
	updateRequest.ID = id
	return &updateRequest, nil
}

type ShowEquipmentRequest struct {
	ID string `json:"id"`
}

func DecodeShowEquipmentRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	return &ShowEquipmentRequest{ID: id}, nil
}

func DecodeDeleteEquipmentRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	return &ShowEquipmentRequest{ID: id}, nil
}

type SupplyCreateRequest struct {
	SupplyID         string `json:"supply_id"`
	SupplyIDCategory string `json:"supply_id_category"`
	Amount           string `json:"amount"`
	SpecialRequest   string `json:"special_request"`
	Zipcode          string `json:"zipcode"`
	CityID           string `json:"city_id"`
	CountyID         string `json:"county_id"`
}

type SupplyUpdateRequest struct {
	ID             string `json:"-"`
	Amount         int    `json:"amount"`
	SpecialRequest string `json:"special_request"`
	Zipcode        string `json:"zipcode"`
}

type ShowSupplyRequest struct {
	ID string
}

func DecodeDeleteSupplyRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	return &ShowSupplyRequest{ID: id}, nil
}

func DecodeShowSupplyRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	return &ShowSupplyRequest{ID: id}, nil
}

func DecodeUpdateSupplyRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	var updateRequest SupplyUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	updateRequest.ID = id
	return &updateRequest, nil
}

func DecodeCreateSupplyRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var create SupplyCreateRequest
	err = json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &create, nil
}

type UserQuotesRequest struct {
	UserID   string
	Accepted bool
	Page     int `json:"page"`
	PerPage  int `json:"per_page"`
}

func DecodeGetUserQuotesRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	accepted, err := strconv.ParseBool(r.URL.Query().Get("accepted"))
	if err != nil {
		accepted = false
	}
	return &UserQuotesRequest{Accepted: accepted,
		Page:    helpers.AlwaysIntFromString(r.URL.Query().Get("page")),
		PerPage: helpers.AlwaysIntFromString(r.URL.Query().Get("per_page"))}, nil
}

type UserQuotesService interface {
	List(ctx context.Context, paginationQuery entities.PaginationQuery, uqRequest UserQuotesRequest) ([]entities.CompleteQuote, error)
	GetPagination(ctx context.Context, uqRequest UserQuotesRequest, paginationQuery entities.PaginationQuery) (*entities.Pagination, error)
}

func MakeUserQuoteEndpoint(svc UserQuotesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ur := request.(*UserQuotesRequest)
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		ur.UserID = userID
		paginationQuery := entities.PaginationQuery{
			Page:    ur.Page,
			PerPage: ur.PerPage,
		}
		quotes, err := svc.List(ctx, paginationQuery, *ur)
		if err != nil {
			return nil, err
		}
		pagination, err := svc.GetPagination(ctx, *ur, paginationQuery)
		if err != nil {
			return nil, err
		}
		quotes = helpers.SortQuotesByStartDate(quotes)
		var i int
		j := pagination.Page * 10
		i = j - 10
		if j > pagination.Total {
			j = pagination.Total
		}
		return responses.CollectionResponse{
			Data:       quotes[i:j],
			Pagination: pagination,
		}, nil
	}
}
