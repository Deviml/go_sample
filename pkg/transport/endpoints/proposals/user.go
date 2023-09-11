package proposals

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/auth/jwt/claims"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
)

type Endpoints struct {
	Create endpoint.Endpoint
	Show   endpoint.Endpoint
	Update endpoint.Endpoint
	Delete endpoint.Endpoint
	Accept endpoint.Endpoint
	Reject endpoint.Endpoint
	List   endpoint.Endpoint
}

type SellerService interface {
	CreateProposal(ctx context.Context, userID string, createRequest CreateRequest) error
	ShowProposal(ctx context.Context, proposalID string) (entities.Proposal, error)
	UpdateProposal(ctx context.Context, updateRequest UpdateRequest) error
	DeleteProposal(ctx context.Context, proposalID string) error
}

type BuyerService interface {
	AcceptProposal(ctx context.Context, proposalID string) error
	RejectProposal(ctx context.Context, proposalID string) error
}

func NewEndpoints(svce SellerService, svcs BuyerService) *Endpoints {
	return &Endpoints{
		Create: MakeCreateProposalRequest(svce),
		Update: MakeUpdateProposalRequest(svce),
		Show:   MakeShowProposalRequest(svce),
		Delete: MakeDeleteProposalRequest(svce),
		Accept: MakeAcceptProposalRequest(svcs),
		Reject: MakeRejectProposalRequest(svcs),
	}
}

func MakeAcceptProposalRequest(svc BuyerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ur := request.(*ShowProposalRequest)
		err = svc.AcceptProposal(ctx, ur.ID)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeRejectProposalRequest(svc BuyerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ur := request.(*ShowProposalRequest)
		err = svc.RejectProposal(ctx, ur.ID)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeDeleteProposalRequest(svc SellerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		dr := request.(*ShowProposalRequest)
		err = svc.DeleteProposal(ctx, dr.ID)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeShowProposalRequest(svc SellerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		sr := request.(*ShowProposalRequest)
		e, err := svc.ShowProposal(ctx, sr.ID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{Data: e}, nil
	}
}

func MakeUpdateProposalRequest(svc SellerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ur := request.(*UpdateRequest)
		err = svc.UpdateProposal(ctx, *ur)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeCreateProposalRequest(svc SellerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		cr := request.(*CreateRequest)
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		err = svc.CreateProposal(ctx, userID, *cr)
		if err != nil {
			return nil, err
		}
		return responses.Created{}, nil
	}
}

type CreateRequest struct {
	QuoteID           string  `json:"quote_id"`
	Make              string  `json:"make"`
	EqModel           string  `json:"eq_model"`
	Year              string  `json:"year"`
	Serial            string  `json:"serial"`
	VIN               string  `json:"vin"`
	EqHours           float32 `json:"eq_hours"`
	Condition         string  `json:"condition"`
	SalePrice         float64 `json:"sale_price"`
	AvailableDate     *int    `json:"available_date"`
	Description       string  `json:"description"`
	Comments          string  `json:"comments"`
	Specifications    string  `json:"specifications"`
	Videos            string  `json:"videos"`
	Pics              string  `json:"pics"`
	Status            string  `json:"status"`
	Freight           string  `json:"freight"`
	Tax               string  `json:"tax"`
	Fees              string  `json:"fees"`
	VendorCompanyName string  `json:"vendor_company_name"`
	VendorEmail       string  `json:"vendor_email"`
	VendorPhoneNumber string  `json:"vendor_phone_number"`
}

func DecodeCreateRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var createRequest CreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		return nil, err
	}
	return &createRequest, nil
}

type UpdateRequest struct {
	ID             string  `json:"-"`
	Make           string  `json:"make"`
	EqModel        string  `json:"eq_model"`
	Year           string  `json:"year"`
	Serial         string  `json:"serial"`
	VIN            string  `json:"vin"`
	EqHours        float32 `json:"eq_hours"`
	Condition      string  `json:"condition"`
	SalePrice      float64 `json:"sale_price"`
	AvailableDate  *int    `json:"available_date"`
	Description    string  `json:"description"`
	Comments       string  `json:"comments"`
	Specifications string  `json:"specifications"`
	Videos         string  `json:"videos"`
	Pics           string  `json:"pics"`
	Freight        string  `json:"freight"`
	Tax            string  `json:"tax"`
	Fees           string  `json:"fees"`
}

func DecodeUpdateRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
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

type ShowProposalRequest struct {
	ID string `json:"id"`
}

func DecodeShowProposalRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	return &ShowProposalRequest{ID: id}, nil
}

func DecodeDeleteRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("bar url request")
	}
	return &ShowProposalRequest{ID: id}, nil
}

type BuyerProposalRequest struct {
	UserID string
	Status string
}

func DecodeGetBuyerProposalRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	status := r.URL.Query().Get("status")
	if err != nil {
		status = ""
	}
	return &BuyerProposalRequest{Status: status}, nil
}

type BuyerProposalService interface {
	Get(ctx context.Context, uqRequest BuyerProposalRequest) ([]entities.Proposal, error)
}

func MakeBuyerProposalEndpoint(svc BuyerProposalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ur := request.(*BuyerProposalRequest)
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		ur.UserID = userID
		proposals, err := svc.Get(ctx, *ur)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{
			Data: proposals,
		}, nil
	}
}

type SellerProposalRequest struct {
	UserID string
	Status string
}

func DecodeGetSellerProposalRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	status := r.URL.Query().Get("status")
	if err != nil {
		status = ""
	}
	return &SellerProposalRequest{Status: status}, nil
}

type SellerProposalService interface {
	Get(ctx context.Context, uqRequest SellerProposalRequest) ([]entities.Proposal, error)
}

func MakeSellerProposalEndpoint(svc SellerProposalService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ur := request.(*SellerProposalRequest)
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		ur.UserID = userID
		proposals, err := svc.Get(ctx, *ur)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{
			Data: proposals,
		}, nil
	}
}
