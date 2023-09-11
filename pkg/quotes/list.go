package quotes

import (
	"context"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/filters"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/sorts"
	mysqlEntities "github.com/Equiphunter-com/equipment-hunter-api/pkg/repositories/mysql/entities"
	"github.com/go-kit/kit/log"
)

type ListQuotesRepository interface {
	List(ctx context.Context, filter filters.ListQuotes, sort sorts.ListQuotes, paginationQuery entities.PaginationQuery) ([]mysqlEntities.Quote, error)
	FetchCount(ctx context.Context, filter filters.ListQuotes) (int, error)
}

type ListService struct {
	logger log.Logger
	lqr    ListQuotesRepository
}

func NewListService(logger log.Logger, lqr ListQuotesRepository) *ListService {
	return &ListService{logger: logger, lqr: lqr}
}

func (l ListService) List(ctx context.Context, filter filters.ListQuotes, sort sorts.ListQuotes, paginationQuery entities.PaginationQuery) ([]entities.Quote, error) {
	repoQuotes, err := l.lqr.List(ctx, filter, sort, paginationQuery)
	if err != nil {
		return nil, err
	}
	return l.buildQuotesFromRepoQuotes(repoQuotes), nil
}

func (l ListService) buildQuotesFromRepoQuotes(repoQuotes []mysqlEntities.Quote) []entities.Quote {
	quotes := make([]entities.Quote, 0, len(repoQuotes))
	for _, repoQuote := range repoQuotes {
		quotes = append(quotes, repoQuote.ToDomainQuote())
	}
	return quotes
}

func (l ListService) GetPagination(ctx context.Context, filter filters.ListQuotes, paginationQuery entities.PaginationQuery) (*entities.Pagination, error) {
	total, err := l.lqr.FetchCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	return entities.MakePaginationFromQueryWithTotal(total, paginationQuery), nil
}
