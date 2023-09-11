package users

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/filters"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/sorts"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/quotes"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/vendorpurchases"
	"github.com/go-kit/kit/log"
)

type SupplyRequestRepository interface {
	GetForUser(ctx context.Context, userID uint, filters *vendorpurchases.GetVendorPurchasesRequest) ([]models.Quote, error)
}

type GetSupplyRequestRepository interface {
	GetByID(ctx context.Context, supplyRequestID string) (*models.Quote, error)
}

type VendorSupplyRequestsService struct {
	logger log.Logger
	srr    SupplyRequestRepository
	gsrr   GetSupplyRequestRepository
	qr     quotes.ListQuotesRepository
}

func NewVendorSupplyRequestsService(logger log.Logger, srr SupplyRequestRepository, gsrr GetSupplyRequestRepository, qr quotes.ListQuotesRepository) *VendorSupplyRequestsService {
	return &VendorSupplyRequestsService{logger: logger, srr: srr, gsrr: gsrr, qr: qr}
}

func (v VendorSupplyRequestsService) GetByID(ctx context.Context, id string) (*entities.CompleteSupplyRequest, error) {
	quote, err := v.gsrr.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	supplyRequest := quote.ToCompleteSupplyRequestEntity()
	return &supplyRequest, nil
}

func (v VendorSupplyRequestsService) GetSupplyRequests(ctx context.Context, request *vendorpurchases.GetVendorPurchasesRequest) ([]entities.CompleteSupplyRequest, error) {
	cq, err := v.qr.List(ctx, filters.ListQuotes{
		Keywords:             request.Keywords,
		Zipcode:              request.Zipcode,
		CityID:               request.CityID,
		EquipmentCategories:  []string{request.EquipmentCategories},
		SupplyCategories:     []string{request.SupplyCategories},
		ContractPreferences:  request.ContractPreferences,
		RentTo:               request.RentTo,
		RentFrom:             request.RentFrom,
		EquipmentSubcategory: request.EquipmentSubcategory,
		UserID:               request.UserID,
		OnlySupply:           true,
		OnlyEquipment:        false,
	}, sorts.ListQuotes{}, entities.PaginationQuery{})
	if err != nil {
		return nil, err
	}
	completeQuotes := make([]entities.CompleteSupplyRequest, 0)

	for _, q := range cq {
		completeQuotes = append(completeQuotes, q.ToCompleteSupplyRequestEntity())
	}

	return completeQuotes, nil
}
