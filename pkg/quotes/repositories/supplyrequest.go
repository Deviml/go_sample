package repositories

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type SupplyRequestRepository struct {
	logger log.Logger
	db     *gorm.DB
}

func NewSupplyRequestRepository(logger log.Logger, db *gorm.DB) *SupplyRequestRepository {
	return &SupplyRequestRepository{logger: logger, db: db}
}

func (s SupplyRequestRepository) GetByID(ctx context.Context, supplyRequestID string) (*models.Quote, error) {
	var quote models.Quote
	result := s.db.
		Preload("WebUser").
		Preload("SupplyRequest.Location").
		Preload("City.State").
		Preload("SupplyRequest.Supply.SupplyCategory").Find(&quote, "id = ?", supplyRequestID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &quote, nil
}
