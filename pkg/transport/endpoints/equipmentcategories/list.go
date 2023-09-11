package equipmentcategories

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/endpoint"
)

type ListService interface {
	List(ctx context.Context) ([]entities.EquipmentCategory, error)
}

func MakeListEndpoint(svc ListService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		equipmentCategories, err := svc.List(ctx)
		if err != nil {
			return nil, err
		}

		return responses.CollectionResponse{
			Data:       equipmentCategories,
			Pagination: nil,
		}, nil
	}
}
