package companies

import (
	"context"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
)

type ListCompanyRepository interface {
	List(ctx context.Context, keyword string) ([]models.Company, error)
}

type ListService struct {
	logger log.Logger
	cr     ListCompanyRepository
}

func NewListService(logger log.Logger, cr ListCompanyRepository) *ListService {
	return &ListService{logger: logger, cr: cr}
}

func (l ListService) List(ctx context.Context, keyword string) ([]entities.Company, error) {
	modelCompanies, err := l.cr.List(ctx, keyword)
	if err != nil {
		return nil, err
	}
	companies := make([]entities.Company, 0)
	for _, modelCompany := range modelCompanies {
		companies = append(companies, modelCompany.ToEntity())
	}
	return companies, nil
}
