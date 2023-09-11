package sublists

import (
	"context"
	"errors"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/filters"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/sorts"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/endpoint"
)

type ListSublistService interface {
	List(ctx context.Context, filters filters.ListSublists, sorts sorts.ListSublists, pagination entities.PaginationQuery) ([]entities.Sublist, error)
	GetPagination(ctx context.Context, filters filters.ListSublists, sorts sorts.ListSublists, pagination entities.PaginationQuery) (*entities.Pagination, error)
}

func MakeListSublistsEndpoint(svc ListSublistService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		getRequest, ok := request.(*requests.GetSublistsRequest)
		if !ok {
			return nil, errors.New("bad cast")
		}
		filters := filters.ListSublists{
			Keywords: getRequest.Keywords,
			Zipcode:  getRequest.Zipcode,
			CityID:   getRequest.CityID,
			StateID:  getRequest.StateID,
		}
		sorts := sorts.ListSublists{
			Criteria:  getRequest.Sort,
			Latitude:  getRequest.Latitude,
			Longitude: getRequest.Longitude,
		}
		paginationQuery := entities.PaginationQuery{
			Page:    getRequest.Page,
			PerPage: getRequest.PerPage,
		}
		sublists, err := svc.List(ctx, filters, sorts, paginationQuery)
		if err != nil {
			return nil, err
		}
		pagination, err := svc.GetPagination(ctx, filters, sorts, paginationQuery)

		if err != nil {
			return nil, err
		}
		return responses.CollectionResponse{
			Data:       makePublishUnix(sublists),
			Pagination: pagination,
		}, nil
	}
}

func makePublishUnix(sublists []entities.Sublist) []entities.Sublist {
	for idx, sublist := range sublists {
		sublist.PublishDateUnix = sublist.PublishDate.Unix()
		sublists[idx] = sublist
	}
	return sublists
}
