package equipmentcategories

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/go-kit/kit/log"
)

type ListEquipmentCategoriesRepository interface {
	List(ctx context.Context) ([]entities.EquipmentCategory, error)
}

type ListService struct {
	logger log.Logger
	lr     ListEquipmentCategoriesRepository
}

func NewListService(logger log.Logger, lr ListEquipmentCategoriesRepository) *ListService {
	return &ListService{logger: logger, lr: lr}
}

func (l ListService) List(ctx context.Context) ([]entities.EquipmentCategory, error) {
	return l.lr.List(ctx)
}
