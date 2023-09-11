package gormclient

import (
	"context"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/vendorpurchases"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type QuotesRepository struct {
	logger log.Logger
	db     *gorm.DB
}

func NewQuotesRepository(logger log.Logger, db *gorm.DB) *QuotesRepository {
	return &QuotesRepository{logger: logger, db: db}
}

//! Get Pagination For Equipment Requests
//? VENDOR DASHBOARD
func (q QuotesRepository) GetPagination(ctx context.Context, eqRequest vendorpurchases.GetVendorPurchasesRequest, paginationQuery entities.PaginationQuery) (*entities.Pagination, error) {
	var user models.WebUser
	result := q.db.Find(&user, "id = ?", eqRequest.UserID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var sc models.Subcontractor
	result = q.db.Find(&sc, "id = ?", user.ProfileID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var quotes []models.Quote
	result = q.db.Preload("SupplyRequest.Location").
		Preload("SupplyRequest.Supply.SupplyCategory").
		Preload("City.State").
		Preload("WebUser").
		Find(&quotes, "web_user_id = ?", eqRequest.UserID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var pd []models.PaymentDetail
	r := q.db.Find(&pd, "id IN ?", quotes)
	if r.Error != nil {
		return nil, r.Error
	}
	for _, v := range pd {
		for i, vv := range quotes {
			if vv.ID == v.ID {
				vv.PurchasedAt = v.CreatedAt
				quotes[i] = vv
			}
		}
	}
	total := len(quotes)
	a := total
	return entities.MakePaginationFromQueryWithTotal(a, paginationQuery), nil
}
func (q QuotesRepository) GetUsersEquipmentRequests(ctx context.Context, userID uint, filters *vendorpurchases.GetVendorPurchasesRequest) ([]models.Quote, error) {
	var wu models.WebUser
	result := q.db.Find(&wu, "id = ?", userID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var user models.VendorRental

	result = q.db.Preload("Quotes", "equipment_request_id IS NOT NULL").
		Preload("Quotes.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory").
		Preload("Quotes.City.State").
		Find(&user, "id = ?", wu.ProfileID)
	if result.Error != nil {
		return nil, result.Error
	}

	qIDs := make([]uint, len(user.Quotes))
	for i, v := range user.Quotes {
		qIDs[i] = v.ID
	}

	var pd []models.PaymentDetail

	r := q.db.Find(&pd, "id IN ?", qIDs)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, v := range pd {
		for i, vv := range user.Quotes {
			if vv.ID == v.ID {
				vv.PurchasedAt = v.CreatedAt
				user.Quotes[i] = vv
			}
		}
	}
	return user.Quotes, nil
}

// Supply
func (q QuotesRepository) GetForUser(ctx context.Context, userID uint, filters *vendorpurchases.GetVendorPurchasesRequest) ([]models.Quote, error) {
	var wu models.WebUser
	result := q.db.Find(&wu, "id = ?", userID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var user models.VendorRental
	result = q.db.Preload("Quotes", "Quotes.supply_request_id IS NOT NULL").
		Preload("Quotes.SupplyRequest.Location").
		Preload("Quotes.City.State").
		Preload("Quotes.SupplyRequest.Supply.SupplyCategory").Find(&user, "id = ?", wu.ProfileID)
	if result.Error != nil {
		return nil, result.Error
	}

	qIDs := make([]uint, len(user.Quotes))
	for i, v := range user.Quotes {
		qIDs[i] = v.ID
	}

	var pd []models.PaymentDetail

	r := q.db.Find(&pd, "id IN ?", qIDs)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, v := range pd {
		for i, vv := range user.Quotes {
			if vv.ID == v.ID {
				vv.PurchasedAt = v.CreatedAt
				user.Quotes[i] = vv
			}
		}
	}

	return user.Quotes, nil
}
