package equipments

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type GetEquipmentService struct {
	logger log.Logger
	db     *gorm.DB
}

func NewGetEquipmentService(logger log.Logger, db *gorm.DB) *GetEquipmentService {
	return &GetEquipmentService{logger: logger, db: db}
}

func (g GetEquipmentService) Get(ctx context.Context, subID string) ([]entities.Equipment, error) {
	var eqs []models.Equipment
	result := g.db
	if subID != "" {
		result = result.Where("equipment_subcategory_id = ?", subID)
	}
	result = result.Find(&eqs)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	eeq := make([]entities.Equipment, 0)
	for _, meq := range eqs {
		eeq = append(eeq, meq.ToEntity())
	}
	return eeq, nil
}
