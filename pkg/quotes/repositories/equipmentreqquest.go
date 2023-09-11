package repositories

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type EquipmentRequestRepository struct {
	logger log.Logger
	db     *gorm.DB
}

func NewEquipmentRequestRepository(logger log.Logger, db *gorm.DB) *EquipmentRequestRepository {
	return &EquipmentRequestRepository{logger: logger, db: db}
}

func (e EquipmentRequestRepository) GetByID(ctx context.Context, id uint) (*models.Quote, error) {
	var quote models.Quote
	result := e.db.
		Preload("WebUser").
		Preload("EquipmentRequest.Location").
		Preload("City.State").
		Preload("EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory").
		Find(&quote, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &quote, nil
}
