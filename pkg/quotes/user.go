package quotes

import (
	"context"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/quotes"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type UserQuotes struct {
	logger log.Logger
	db     *gorm.DB
}

func NewUserQuotes(logger log.Logger, db *gorm.DB) *UserQuotes {
	return &UserQuotes{logger: logger, db: db}
}

func (u UserQuotes) List(ctx context.Context, paginationQuery entities.PaginationQuery, uqRequest quotes.UserQuotesRequest) ([]entities.CompleteQuote, error) {
	var user models.WebUser
	result := u.db.Find(&user, "id = ?", uqRequest.UserID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var sc models.Subcontractor
	result = u.db.Find(&sc, "id = ?", user.ProfileID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var quotes []models.Quote
	result = u.db.Preload("EquipmentRequest.Location").
		Preload("EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory").
		Preload("SupplyRequest.Location").
		Preload("SupplyRequest.Supply.SupplyCategory").
		Preload("VendorRentals").
		Preload("City.State").
		Preload("WebUser").
		Find(&quotes, "web_user_id = ?", uqRequest.UserID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	eQ := make([]entities.CompleteQuote, 0)
	for _, mq := range quotes {
		if uqRequest.Accepted && len(mq.VendorRentals) == 0 {
			continue
		}
		if !uqRequest.Accepted && len(mq.VendorRentals) > 0 {
			continue
		}
		eQ = append(eQ, mq.ToQuiteEntity())
	}

	return eQ, nil
}

func (u UserQuotes) GetPagination(ctx context.Context, uqRequest quotes.UserQuotesRequest, paginationQuery entities.PaginationQuery) (*entities.Pagination, error) {
	var user models.WebUser
	result := u.db.Find(&user, "id = ?", uqRequest.UserID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var sc models.Subcontractor
	result = u.db.Find(&sc, "id = ?", user.ProfileID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var quotes []models.Quote
	result = u.db.Preload("EquipmentRequest.Location").
		Preload("EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory").
		Preload("SupplyRequest.Location").
		Preload("SupplyRequest.Supply.SupplyCategory").
		Preload("VendorRentals").
		Preload("City.State").
		Preload("WebUser").
		Find(&quotes, "web_user_id = ?", uqRequest.UserID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	eQ := make([]entities.CompleteQuote, 0)
	for _, mq := range quotes {
		if uqRequest.Accepted && len(mq.VendorRentals) == 0 {
			continue
		}
		if !uqRequest.Accepted && len(mq.VendorRentals) > 0 {
			continue
		}
		eQ = append(eQ, mq.ToQuiteEntity())
	}
	total := len(eQ)
	return entities.MakePaginationFromQueryWithTotal(total, paginationQuery), nil
}
