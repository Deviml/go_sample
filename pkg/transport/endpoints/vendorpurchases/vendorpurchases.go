package vendorpurchases

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

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
	GetSublists                     endpoint.Endpoint
	GetEquipmentRequestsSortedRD    endpoint.Endpoint
	GetEquipmentRequestsSortedRA    endpoint.Endpoint
	GetEquipmentRequestsSortedPA    endpoint.Endpoint
	GetEquipmentRequestsSortedPD    endpoint.Endpoint
	GetSupplyRequests               endpoint.Endpoint
	GetSublistsDetail               endpoint.Endpoint
	GetSupplyDetail                 endpoint.Endpoint
	GetEquipmentDetail              endpoint.Endpoint
	Checkout                        endpoint.Endpoint
	PaymentIntent                   endpoint.Endpoint
	ConfirmPaymentIntentAndCheckout endpoint.Endpoint
}

type CheckoutService interface {
	Checkout(ctx context.Context, userID string, request CheckoutRequest) error
	PaymentIntent(ctx context.Context, userID string, request PaymentIntentRequest) (interface{}, error)
	ConfirmPaymentIntentAndCheckout(ctx context.Context, userID string, request CheckoutRequest) error
}

func NewEndpoints(ssvc GetVendorSublistsService, ersvc GetEquipmentRequestsService, srsvc GetSupplyRequestsService, csvc CheckoutService) *Endpoints {
	return &Endpoints{
		GetSublists:                     MakeGetSublistsEndpoint(ssvc),
		GetSupplyRequests:               MakeGetSupplyRequestEndpoint(srsvc),
		GetEquipmentRequestsSortedRD:    MakeEquipmentRequestByRequestDateDescEndpoint(ersvc),
		GetEquipmentRequestsSortedRA:    MakeEquipmentRequestByRequestDateAscEndpoint(ersvc),
		GetEquipmentRequestsSortedPA:    MakeEquipmentRequestByPurchaseDateAscEndpoint(ersvc),
		GetEquipmentRequestsSortedPD:    MakeEquipmentRequestByPurchaseDateDescEndpoint(ersvc),
		GetSublistsDetail:               MakeGetSublistsDetailEndpoint(ssvc),
		GetSupplyDetail:                 MakeGetSupplyDetailEndpoint(srsvc),
		GetEquipmentDetail:              MakeEquipmentDetailEndpoint(ersvc),
		Checkout:                        MakeCheckoutEndpoint(csvc),
		PaymentIntent:                   MakePaymentIntentEndpoint(csvc),
		ConfirmPaymentIntentAndCheckout: MakeConfirmPaymentIntentAndCheckoutEndpoint(csvc),
	}
}

type CheckoutRequest struct {
	Email           string                 `json:"email"`
	NameOnCard      string                 `json:"name_on_card"`
	Zipcode         string                 `json:"zipcode"`
	Quotes          []string               `json:"quotes"`
	Sublists        []int                  `json:"sublists"`
	Total           float32                `json:"total"`
	Token           map[string]interface{} `json:"token"`
	PaymentResponse map[string]interface{} `json:"paymentResponse"`
}

func MakeCheckoutEndpoint(svc CheckoutService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cr := request.(*CheckoutRequest)
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		err = svc.Checkout(ctx, userID, *cr)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeConfirmPaymentIntentAndCheckoutEndpoint(svc CheckoutService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cr := request.(*CheckoutRequest)
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		err = svc.ConfirmPaymentIntentAndCheckout(ctx, userID, *cr)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

type PaymentIntentRequest struct {
	Quotes   []string `json:"quotes"`
	Sublists []int    `json:"sublists"`
}

func MakePaymentIntentEndpoint(svc CheckoutService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		pir := request.(*PaymentIntentRequest)
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token == nil {
			return nil, errors.New("token is not present")
		}
		userToken := token.(*claims.UserClaims)
		userID = userToken.UserID
		paymentIntent, err := svc.PaymentIntent(ctx, userID, *pir)

		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{Data: paymentIntent}, nil
	}
}

type GetVendorPurchasesRequest struct {
	Keywords             []string `json:"keywords"`
	From                 int64    `json:"from"`
	To                   int64    `json:"to"`
	UserID               string   `json:"user_id"`
	Zipcode              string   `json:"zipcode"`
	CityID               int      `json:"city"`
	Sort                 string   `json:"sort"`
	Latitude             float32  `json:"lat"`
	Longitude            float32  `json:"lon"`
	Page                 int      `json:"page"`
	PerPage              int      `json:"per_page"`
	EquipmentCategories  string   `json:"equipment_categories"`
	SupplyCategories     string   `json:"supply_categories"`
	ContractPreferences  []int    `json:"contract_preferences"`
	RentTo               int64    `json:"rent_to"`
	RentFrom             int64    `json:"rent_from"`
	EquipmentSubcategory string   `json:"equipment_subcategory"`
}

func DecodeGetVendorSublistsRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	contractPreferences := make([]int, 0)
	filterContractPreferences, ok := r.URL.Query()["contract_preference[]"]
	if ok {
		for _, filterContractPreference := range filterContractPreferences {
			contractPreferences = append(contractPreferences, helpers.AlwaysIntFromString(filterContractPreference))
		}
	}
	return &GetVendorPurchasesRequest{
		From:                 helpers.AlwaysInt64FromString(r.URL.Query().Get("from")),
		To:                   helpers.AlwaysInt64FromString(r.URL.Query().Get("to")),
		RentFrom:             helpers.AlwaysInt64FromString(r.URL.Query().Get("rent_from")),
		RentTo:               helpers.AlwaysInt64FromString(r.URL.Query().Get("rent_to")),
		Keywords:             helpers.SplitStringQuery(r.URL.Query().Get("keywords")),
		Zipcode:              r.URL.Query().Get("zipcode"),
		CityID:               helpers.AlwaysIntFromString(r.URL.Query().Get("city")),
		Latitude:             helpers.AlwaysFloatFromString(r.URL.Query().Get("lat")),
		Longitude:            helpers.AlwaysFloatFromString(r.URL.Query().Get("lon")),
		Page:                 helpers.AlwaysIntFromString(r.URL.Query().Get("page")),
		PerPage:              helpers.AlwaysIntFromString(r.URL.Query().Get("per_page")),
		EquipmentCategories:  r.URL.Query().Get("equipment_categories"),
		SupplyCategories:     r.URL.Query().Get("supply_categories"),
		ContractPreferences:  contractPreferences,
		EquipmentSubcategory: r.URL.Query().Get("equipment_subcategory"),
	}, nil
}

type PurchaseDetailRequest struct {
	ID string `json:"id"`
}

func DecodePurchaseDetailRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	return &PurchaseDetailRequest{ID: id}, nil
}

func DecodeCheckoutRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var c CheckoutRequest
	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &c, nil
}

func DecodePaymentIntentRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var p PaymentIntentRequest
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &p, nil
}

type GetVendorSublistsService interface {
	GetSublists(ctx context.Context, request *GetVendorPurchasesRequest) ([]entities.Sublist, error)
	GetByID(ctx context.Context, id string) (*entities.CompleteSublist, error)
}

func MakeGetSublistsEndpoint(svc GetVendorSublistsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getVendorSublistsRequest, ok := request.(*GetVendorPurchasesRequest)
		if !ok {
			return nil, http2.BadRequest{}
		}
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		getVendorSublistsRequest.UserID = userID
		sublists, err := svc.GetSublists(ctx, getVendorSublistsRequest)
		if err != nil {
			return nil, err
		}
		return responses.CollectionResponse{
			Data: sublists,
			Pagination: &entities.Pagination{
				Page:       1,
				PerPage:    len(sublists),
				TotalPages: 1,
				Total:      len(sublists),
			},
		}, nil
	}
}

type GetSupplyRequestsService interface {
	GetSupplyRequests(ctx context.Context, request *GetVendorPurchasesRequest) ([]entities.CompleteSupplyRequest, error)
	GetByID(ctx context.Context, id string) (*entities.CompleteSupplyRequest, error)
}

func MakeGetSupplyRequestEndpoint(svc GetSupplyRequestsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getVendorSublistsRequest, ok := request.(*GetVendorPurchasesRequest)
		if !ok {
			return nil, http2.BadRequest{}
		}
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		getVendorSublistsRequest.UserID = userID
		sublists, err := svc.GetSupplyRequests(ctx, getVendorSublistsRequest)
		if err != nil {
			return nil, err
		}
		return responses.CollectionResponse{
			Data: sublists,
			Pagination: &entities.Pagination{
				Page:       1,
				PerPage:    len(sublists),
				TotalPages: 1,
				Total:      len(sublists),
			},
		}, nil
	}
}

type GetEquipmentRequestsService interface {
	GetEquipmentRequests(ctx context.Context, request GetVendorPurchasesRequest) ([]entities.CompleteEquipmentRequest, error)
	GetByID(ctx context.Context, id string) (*entities.CompleteEquipmentRequest, error)
	GetPagination(ctx context.Context, request GetVendorPurchasesRequest, paginationQuery entities.PaginationQuery) (*entities.Pagination, error)
}

// TODO: Working here
func MakeEquipmentRequestByPurchaseDateAscEndpoint(svc GetEquipmentRequestsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getVendorEquipmentRequest, ok := request.(*GetVendorPurchasesRequest)
		if !ok {
			return nil, http2.BadRequest{}
		}
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		getVendorEquipmentRequest.UserID = userID
		paginationQuery := entities.PaginationQuery{
			Page:    getVendorEquipmentRequest.Page,
			PerPage: getVendorEquipmentRequest.PerPage,
		}
		eRequests, err := svc.GetEquipmentRequests(ctx, *getVendorEquipmentRequest)
		if err != nil {
			return nil, err
		}
		pagination, err := svc.GetPagination(ctx, *getVendorEquipmentRequest, paginationQuery)
		if err != nil {
			return nil, err
		}
		pagination.Total = len(eRequests)
		eRequests = helpers.SortRequestsByPurchaseDateAsc(eRequests)
		var i int
		j := pagination.Page * 10
		i = j - 10
		if j > pagination.Total {
			j = pagination.Total
		}
		return responses.CollectionResponse{
			Data:       eRequests[i:j],
			Pagination: pagination,
		}, nil
	}
}

func MakeEquipmentRequestByRequestDateAscEndpoint(svc GetEquipmentRequestsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getVendorEquipmentRequest, ok := request.(*GetVendorPurchasesRequest)
		if !ok {
			return nil, http2.BadRequest{}
		}
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		getVendorEquipmentRequest.UserID = userID
		paginationQuery := entities.PaginationQuery{
			Page:    getVendorEquipmentRequest.Page,
			PerPage: getVendorEquipmentRequest.PerPage,
		}
		eRequests, err := svc.GetEquipmentRequests(ctx, *getVendorEquipmentRequest)
		if err != nil {
			return nil, err
		}
		pagination, err := svc.GetPagination(ctx, *getVendorEquipmentRequest, paginationQuery)
		if err != nil {
			return nil, err
		}
		pagination.Total = len(eRequests)
		eRequests = helpers.SortRequestsByRequestDateAsc(eRequests)
		var i int
		j := pagination.Page * 10
		i = j - 10
		if j > pagination.Total {
			j = pagination.Total
		}
		return responses.CollectionResponse{
			Data:       eRequests[i:j],
			Pagination: pagination,
		}, nil
	}
}

func MakeEquipmentRequestByRequestDateDescEndpoint(svc GetEquipmentRequestsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getVendorEquipmentRequest, ok := request.(*GetVendorPurchasesRequest)
		if !ok {
			return nil, http2.BadRequest{}
		}
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		getVendorEquipmentRequest.UserID = userID
		paginationQuery := entities.PaginationQuery{
			Page:    getVendorEquipmentRequest.Page,
			PerPage: getVendorEquipmentRequest.PerPage,
		}
		eRequests, err := svc.GetEquipmentRequests(ctx, *getVendorEquipmentRequest)
		if err != nil {
			return nil, err
		}
		pagination, err := svc.GetPagination(ctx, *getVendorEquipmentRequest, paginationQuery)
		if err != nil {
			return nil, err
		}
		pagination.Total = len(eRequests)
		eRequests = helpers.SortRequestsByRequestDateDesc(eRequests)
		var i int
		j := pagination.Page * 10
		i = j - 10
		if j > pagination.Total {
			j = pagination.Total
		}
		return responses.CollectionResponse{
			Data:       eRequests[i:j],
			Pagination: pagination,
		}, nil
	}
}

// Pagination For Vendor Dashboard
func MakeEquipmentRequestByPurchaseDateDescEndpoint(svc GetEquipmentRequestsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getVendorEquipmentRequest, ok := request.(*GetVendorPurchasesRequest)
		if !ok {
			return nil, http2.BadRequest{}
		}
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		getVendorEquipmentRequest.UserID = userID
		paginationQuery := entities.PaginationQuery{
			Page:    getVendorEquipmentRequest.Page,
			PerPage: getVendorEquipmentRequest.PerPage,
		}
		eRequests, err := svc.GetEquipmentRequests(ctx, *getVendorEquipmentRequest)
		if err != nil {
			return nil, err
		}
		pagination, err := svc.GetPagination(ctx, *getVendorEquipmentRequest, paginationQuery)
		if err != nil {
			return nil, err
		}
		pagination.Total = len(eRequests)
		eRequests = helpers.SortRequestsByPurchaseDateDesc(eRequests)
		var i int
		j := pagination.Page * 10
		i = j - 10
		if j > pagination.Total {
			j = pagination.Total
		}
		return responses.CollectionResponse{
			Data:       eRequests[i:j],
			Pagination: pagination,
		}, nil
	}
}

func MakeGetSublistsDetailEndpoint(svc GetVendorSublistsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getByIDRequest := request.(*PurchaseDetailRequest)
		sublist, err := svc.GetByID(ctx, getByIDRequest.ID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{
			Data: sublist,
		}, nil
	}
}

func MakeGetSupplyDetailEndpoint(svc GetSupplyRequestsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getByIDRequest := request.(*PurchaseDetailRequest)
		sublist, err := svc.GetByID(ctx, getByIDRequest.ID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{
			Data: sublist,
		}, nil
	}
}

func MakeEquipmentDetailEndpoint(svc GetEquipmentRequestsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getByIDRequest := request.(*PurchaseDetailRequest)
		sublist, err := svc.GetByID(ctx, getByIDRequest.ID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{
			Data: sublist,
		}, nil
	}
}
