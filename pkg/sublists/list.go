package sublists

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/filters"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/sorts"
	"github.com/go-kit/kit/log"
)

type ListSublistRepository interface {
	FetchSublists(ctx context.Context, filters filters.ListSublists, sorts sorts.ListSublists, pagination entities.PaginationQuery) ([]entities.Sublist, error)
	FetchCount(ctx context.Context, filters filters.ListSublists, sorts sorts.ListSublists) (int, error)
}

type NopListSublistRepository struct {
}

func (n NopListSublistRepository) FetchSublists(ctx context.Context, filters filters.ListSublists, sorts sorts.ListSublists, pagination entities.PaginationQuery) ([]entities.Sublist, error) {
	return []entities.Sublist{}, nil
}

func (n NopListSublistRepository) FetchCount(ctx context.Context, filters filters.ListSublists, sorts sorts.ListSublists) (int, error) {
	return 0, nil
}

type ListSublistsService struct {
	logger log.Logger
	lsr    ListSublistRepository
}

func NewListSublistsService(logger log.Logger, lsr ListSublistRepository) *ListSublistsService {
	if lsr == nil {
		lsr = NopListSublistRepository{}
	}
	return &ListSublistsService{logger: logger, lsr: lsr}
}

func (l ListSublistsService) GetPagination(ctx context.Context, filters filters.ListSublists, sorts sorts.ListSublists, pagination entities.PaginationQuery) (*entities.Pagination, error) {
	total, err := l.lsr.FetchCount(ctx, filters, sorts)
	if err != nil {
		return nil, err
	}
	return entities.MakePaginationFromQueryWithTotal(total, pagination), nil
}

func (l ListSublistsService) List(ctx context.Context, filters filters.ListSublists, sorts sorts.ListSublists, pagination entities.PaginationQuery) ([]entities.Sublist, error) {
	return l.lsr.FetchSublists(ctx, filters, sorts, pagination)
}
