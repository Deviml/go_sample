package quotes

import (
	"context"
	"errors"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/filters"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/sorts"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/helpers"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/endpoint"
)

type ListQuotesService interface {
	List(ctx context.Context, filter filters.ListQuotes, sort sorts.ListQuotes, paginationQuery entities.PaginationQuery) ([]entities.Quote, error)
	GetPagination(ctx context.Context, filter filters.ListQuotes, paginationQuery entities.PaginationQuery) (*entities.Pagination, error)
}

func MakeListQuotesEndpoint(svc ListQuotesService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		listRequest, ok := request.(*requests.ListQuotesRequest)
		if !ok {
			return nil, errors.New("bad cast")
		}
		filter := filters.MakeListQuotesFromRequest(*listRequest)
		sort := sorts.MakeListQuotesFromRequest(*listRequest)

		paginationQuery := entities.PaginationQuery{
			Page:    listRequest.Page,
			PerPage: listRequest.PerPage,
		}

		newpaginationQuery := entities.PaginationQuery{
			Page:    1,
			PerPage: 201,
		}

		// quotes, err := svc.List(ctx, filter, sort, paginationQuery)
		// if err != nil {
		// 	return nil, err
		// }

		newquotes, _ := svc.List(ctx, filter, sort, newpaginationQuery)
		newquotes = helpers.UpdateListQuotesExpiryStatus(newquotes)

		if sort.Criteria == "name" {

			newquotes = helpers.SortQuotesByNameAsc(newquotes)
		} else if sort.Criteria == "date" {

			newquotes = helpers.SortQuotesByExpirationDateDesc(newquotes)
		} else if sort.Criteria == "location" {

			newquotes = helpers.SortQuotesByLocationAsc(newquotes)
		}

		pagination, err := svc.GetPagination(ctx, filter, paginationQuery)
		if err != nil {
			return nil, err
		}
		pagination.Total = len(newquotes)
		var i int
		j := pagination.Page * 10
		i = j - 10
		if j > pagination.Total {
			j = pagination.Total
		}
		// quotes = helpers.UpdateListQuotesExpiryStatus(quotes)
		return responses.CollectionResponse{
			Data:       newquotes[i:j],
			Pagination: pagination,
		}, nil
	}
}
