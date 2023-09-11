package supplycategories

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)

type ListSupplyCategoriesService interface {
	ListSupplyCategories(ctx context.Context) ([]entities.SupplyCategory, error)
	GetSupplies(ctx context.Context, sID string) ([]entities.Supply, error)
}

func MakeListSupplyCategoriesEndpoint(svc ListSupplyCategoriesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		supplyCategories, err := svc.ListSupplyCategories(ctx)
		if err != nil {
			return nil, err
		}
		return responses.CollectionResponse{
			Data:       supplyCategories,
			Pagination: nil,
		}, nil
	}
}

type GetSupplies struct {
	SubcategoryID string
}

func DecodeSupplies(ctx context.Context, r *http.Request) (response interface{}, err error) {
	return &GetSupplies{
		SubcategoryID: r.URL.Query().Get("subcategory"),
	}, nil
}

func MakeSuppliesEndpoint(svc ListSupplyCategoriesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*GetSupplies)
		s, err := svc.GetSupplies(ctx, r.SubcategoryID)
		if err != nil {
			return nil, err
		}
		return responses.CollectionResponse{
			Data:       s,
			Pagination: nil,
		}, nil
	}
}
