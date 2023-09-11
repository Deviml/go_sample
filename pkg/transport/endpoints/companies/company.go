package companies

import (
	"context"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/endpoint"
)

type Services struct {
	ListCompaniesService
}

func NewServices(listCompaniesService ListCompaniesService) *Services {
	return &Services{ListCompaniesService: listCompaniesService}
}

type ListCompaniesService interface {
	List(ctx context.Context, keyword string) ([]entities.Company, error)
}

type Endpoints struct {
	List endpoint.Endpoint
}

func NewEndpoints(svc *Services) *Endpoints {
	return &Endpoints{List: MakeListCompaniesEndpoint(svc)}
}

func MakeListCompaniesEndpoint(svc ListCompaniesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		listRequest := request.(*requests.ListCompaniesRequest)
		companies, err := svc.List(ctx, listRequest.Keyword)
		if err != nil {
			return nil, err
		}
		return responses.CollectionResponse{
			Data: companies,
		}, nil
	}
}
