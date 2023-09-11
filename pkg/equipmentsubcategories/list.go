package equipmentsubcategories

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/go-kit/kit/log"
)

type GetEquipmentSubCategory interface {
	GetList(ctx context.Context, equipmentCategoryID string) ([]entities.EquipmentSubcategory, error)
}

type GetService struct {
	logger log.Logger
	gesbr  GetEquipmentSubCategory
}

func NewGetService(logger log.Logger, gesbr GetEquipmentSubCategory) *GetService {
	return &GetService{logger: logger, gesbr: gesbr}
}

func (g GetService) GetEquipmentSubcategories(ctx context.Context, equipmentCategory string) ([]entities.EquipmentSubcategory, error) {
	return g.gesbr.GetList(ctx, equipmentCategory)
}
