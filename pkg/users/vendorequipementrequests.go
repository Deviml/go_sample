package users

import (
	"context"
	"strconv"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/filters"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/sorts"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/quotes"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/vendorpurchases"
	"github.com/go-kit/kit/log"
)

type EquipmentRequestRepository interface {
	GetUsersEquipmentRequests(ctx context.Context, userID uint, filters *vendorpurchases.GetVendorPurchasesRequest) ([]models.Quote, error)
	GetPagination(ctx context.Context, eqRequest vendorpurchases.GetVendorPurchasesRequest, paginationQuery entities.PaginationQuery) (*entities.Pagination, error)
}

type GetEquipmentRequestRepository interface {
	GetByID(ctx context.Context, id uint) (*models.Quote, error)
}

type EquipmentRequestsService struct {
	err    EquipmentRequestRepository
	logger log.Logger
	gerr   GetEquipmentRequestRepository
	qr     quotes.ListQuotesRepository
}

func NewEquipmentRequestsService(err EquipmentRequestRepository, logger log.Logger, gerr GetEquipmentRequestRepository, qr quotes.ListQuotesRepository) *EquipmentRequestsService {
	return &EquipmentRequestsService{err: err, logger: logger, gerr: gerr, qr: qr}
}

func (e EquipmentRequestsService) GetPagination(ctx context.Context, eqRequest vendorpurchases.GetVendorPurchasesRequest, paginationQuery entities.PaginationQuery) (*entities.Pagination, error) {
	cq, err := e.qr.List(ctx, filters.ListQuotes{
		Keywords:             eqRequest.Keywords,
		Zipcode:              eqRequest.Zipcode,
		CityID:               eqRequest.CityID,
		EquipmentCategories:  []string{eqRequest.EquipmentCategories},
		SupplyCategories:     []string{eqRequest.SupplyCategories},
		ContractPreferences:  eqRequest.ContractPreferences,
		RentTo:               eqRequest.RentTo,
		RentFrom:             eqRequest.RentFrom,
		EquipmentSubcategory: eqRequest.EquipmentSubcategory,
		UserID:               eqRequest.UserID,
		OnlySupply:           false,
		OnlyEquipment:        true,
		NotExpired:           true,
		NotServed:            true,
	}, sorts.ListQuotes{}, paginationQuery)
	if err != nil {
		return nil, err
	}
	completeQuotes := make([]entities.CompleteEquipmentRequest, 0)

	for _, q := range cq {
		completeQuotes = append(completeQuotes, q.ToCompleteEquipmentRequestEntity())
	}

	total := len(completeQuotes)
	return entities.MakePaginationFromQueryWithTotal(total, paginationQuery), nil
}

func (e EquipmentRequestsService) GetByID(ctx context.Context, id string) (*entities.CompleteEquipmentRequest, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	quote, err := e.gerr.GetByID(ctx, uint(uid))
	if err != nil {
		return nil, err
	}
	er := quote.ToCompleteEquipmentRequestEntity()
	return &er, nil
}

func (e EquipmentRequestsService) GetEquipmentRequests(ctx context.Context, request vendorpurchases.GetVendorPurchasesRequest) ([]entities.CompleteEquipmentRequest, error) {
	cq, err := e.qr.List(ctx, filters.ListQuotes{
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
		OnlySupply:           false,
		OnlyEquipment:        true,
		NotExpired:           true,
		NotServed:            true,
	}, sorts.ListQuotes{}, entities.PaginationQuery{})
	if err != nil {
		return nil, err
	}
	completeQuotes := make([]entities.CompleteEquipmentRequest, 0)

	for _, q := range cq {
		completeQuotes = append(completeQuotes, q.ToCompleteEquipmentRequestEntity())
	}

	return completeQuotes, nil
}
