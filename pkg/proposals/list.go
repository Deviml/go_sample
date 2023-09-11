package proposals

import (
	"context"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	mysqlEntities "github.com/Equiphunter-com/equipment-hunter-api/pkg/repositories/mysql/entities"
	"github.com/go-kit/kit/log"
)

type ListProposalsRepository interface {
	List(ctx context.Context, userID string, paginationQuery entities.PaginationQuery, userType string) ([]mysqlEntities.Proposal, error)
	FetchCount(ctx context.Context, userID string, userType string) (int, error)
}

type ListService struct {
	logger log.Logger
	lqr    ListProposalsRepository
}

func NewListService(logger log.Logger, lqr ListProposalsRepository) *ListService {
	return &ListService{logger: logger, lqr: lqr}
}

func (l ListService) List(ctx context.Context, userID string, paginationQuery entities.PaginationQuery, userType string) ([]entities.Proposal, error) {
	repoProposals, err := l.lqr.List(ctx, userID, paginationQuery, userType)
	if err != nil {
		return nil, err
	}
	return l.buildProposalsFromRepoProposals(repoProposals), nil
}

func (l ListService) buildProposalsFromRepoProposals(repoProposals []mysqlEntities.Proposal) []entities.Proposal {
	proposals := make([]entities.Proposal, 0, len(repoProposals))
	for _, repoProposal := range repoProposals {
		proposals = append(proposals, repoProposal.ToListProposals())
	}
	return proposals
}

func (l ListService) GetPagination(ctx context.Context, userID string, paginationQuery entities.PaginationQuery, userType string) (*entities.Pagination, error) {
	total, err := l.lqr.FetchCount(ctx, userID, userType)
	if err != nil {
		return nil, err
	}
	return entities.MakePaginationFromQueryWithTotal(total, paginationQuery), nil
}
