package repositories

import (
	"context"
	"fmt"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type ListCompaniesRepository struct {
	logger log.Logger
	db     *gorm.DB
}

func NewListCompaniesRepository(logger log.Logger, db *gorm.DB) *ListCompaniesRepository {
	return &ListCompaniesRepository{logger: logger, db: db}
}

func (l ListCompaniesRepository) List(ctx context.Context, keyword string) ([]models.Company, error) {
	var companies []models.Company
	query := l.db
	if keyword != "" {
		query = query.Where("company_name LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	result := query.Find(&companies)
	if result.Error != nil {
		return nil, result.Error
	}
	return companies, nil
}
