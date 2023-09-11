package repositories

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"gorm.io/gorm"
)

type SingleQuoteRepository struct {
	logger log.Logger
	db     *gorm.DB
}

func NewSingleQuoteRepository(logger log.Logger, db *gorm.DB) *SingleQuoteRepository {
	return &SingleQuoteRepository{logger: logger, db: db}
}

func (s SingleQuoteRepository) GetByID(ctx context.Context, id string) (entities.CompleteSingleQuote, error) {

	var quote models.Quote

	result := s.db.Preload("EquipmentRequest.Location").Preload("EquipmentRequest.Equipment.EquipmentSubcategory").Preload("EquipmentRequest.Equipment").Preload("City.State").Find(&quote, "id = ?", id)

	if result != nil && result.Error != nil {
		return entities.CompleteSingleQuote{}, result.Error
	}
	return quote.ToSingleQuoteEntity(), nil
}
