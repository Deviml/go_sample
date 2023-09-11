package quotes

import (
	"context"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/endpoint"
)

type Services struct {
	GetByIDService
}

func NewServices(getByIDService GetByIDService) *Services {
	return &Services{GetByIDService: getByIDService}
}

type GetByIDService interface {
	GetByID(ctx context.Context, id string) (entities.CompleteSingleQuote, error)
}

type Endpoint struct {
	GetByID endpoint.Endpoint
}

func NewSingleEndpoints(svc *Services) *Endpoint {
	return &Endpoint{GetByID: MakeGetByIDEndpoint(svc)}
}

func MakeGetByIDEndpoint(svc GetByIDService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getByIDRequest := request.(*requests.GetByIDRequest)
		quote, err := svc.GetByID(ctx, getByIDRequest.ID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{
			Data: quote,
		}, nil
	}
}
