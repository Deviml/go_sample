package equipmentsubcategories

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/endpoint"
)

type GetEquipmentSubcategoriesService interface {
	GetEquipmentSubcategories(ctx context.Context, equipmentCategory string) ([]entities.EquipmentSubcategory, error)
}

func MakeGetEquipmentSubcategoriesEndpoint(svc GetEquipmentSubcategoriesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getRequest, _ := request.(*requests.GetEquipmentSubcategories)
		equipmentSubcategories, err := svc.GetEquipmentSubcategories(ctx, getRequest.CategoryID)
		if err != nil {
			return nil, err
		}
		return &responses.CollectionResponse{
			Data:       equipmentSubcategories,
			Pagination: nil,
		}, nil
	}
}
