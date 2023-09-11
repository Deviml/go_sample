package cities

import (
	"context"
	http2 "net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/helpers"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetCities   endpoint.Endpoint
	GetCounties endpoint.Endpoint
	GetState    endpoint.Endpoint
}

func NewEndpoints(getCities endpoint.Endpoint, getCounties endpoint.Endpoint, getStates endpoint.Endpoint) *Endpoints {
	return &Endpoints{GetCities: getCities, GetCounties: getCounties, GetState: getStates}
}

type GetCitiesService interface {
	ListCities(ctx context.Context, keywords []string, stateID string) ([]entities.City, error)
	ListCounties(ctx context.Context, keywords []string, stateID string) ([]entities.County, error)
	ListState(ctx context.Context, keywords []string) ([]entities.State, error)
}

type GetServicesRequest struct {
	Keywords []string `json:"keywords"`
	StateID  string
}

func DecodeGetCitiesRequest(ctx context.Context, r *http2.Request) (response interface{}, err error) {
	return &GetServicesRequest{
		Keywords: helpers.SplitStringQuery(r.URL.Query().Get("keywords")),
		StateID:  r.URL.Query().Get("state_id"),
	}, nil
}

func MakeGetCitiesEndpoint(svc GetCitiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r, ok := request.(*GetServicesRequest)
		if !ok {
			return nil, err
		}
		cities, err := svc.ListCities(ctx, r.Keywords, r.StateID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{Data: cities}, nil
	}
}

func MakeGetCountiesEndpoint(svc GetCitiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r, ok := request.(*GetServicesRequest)
		if !ok {
			return nil, err
		}
		cities, err := svc.ListCounties(ctx, r.Keywords, r.StateID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{Data: cities}, nil
	}
}

func MakeGetStatesEndpoint(svc GetCitiesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r, ok := request.(*GetServicesRequest)
		if !ok {
			return nil, err
		}
		cities, err := svc.ListState(ctx, r.Keywords)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{Data: cities}, nil
	}
}
