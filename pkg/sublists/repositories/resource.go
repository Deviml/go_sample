package repositories

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type SublistRepository struct {
	logger log.Logger
	db     *gorm.DB
}

func NewSublistRepository(logger log.Logger, db *gorm.DB) *SublistRepository {
	return &SublistRepository{logger: logger, db: db}
}

func (s SublistRepository) GetSublistByID(ctx context.Context, sublistID string) (*models.Sublist, error) {
	var sublist models.Sublist
	result := s.db.Preload("Location").Preload("City.State").Preload("WebUser").Find(&sublist, "id = ?", sublistID)
	if result.Error != nil {
		return nil, result.Error
	}
	var sc []models.SublistsCompany
	result = s.db.Preload("Company").Find(&sc, "sublist_id = ?", sublistID)
	if result.Error != nil {
		return nil, result.Error
	}
	sublist.SublistsCompanies = sc
	return &sublist, nil
}
