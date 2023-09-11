package supplycategories

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type ListSupplyCategoriesRepositories interface {
	List(ctx context.Context) ([]entities.SupplyCategory, error)
}

type ListService struct {
	logger log.Logger
	lr     ListSupplyCategoriesRepositories
	db     *gorm.DB
}

func NewListService(logger log.Logger, lr ListSupplyCategoriesRepositories, db *gorm.DB) *ListService {
	return &ListService{logger: logger, lr: lr, db: db}
}

func (l ListService) GetSupplies(ctx context.Context, sID string) ([]entities.Supply, error) {
	var ms []models.Supply
	r := l.db
	if sID != "" {
		r = r.Where("supply_category_id = ?", sID)
	}
	result := r.Find(&ms)
	if result.Error != nil {
		return nil, result.Error
	}
	s := make([]entities.Supply, 0)
	for _, m := range ms {
		s = append(s, m.ToEntity())
	}
	return s, nil
}

func (l ListService) ListSupplyCategories(ctx context.Context) ([]entities.SupplyCategory, error) {
	return l.lr.List(ctx)
}
