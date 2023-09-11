package proposals

import (
	"context"
	"database/sql"
	"math/rand"
	"strconv"
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/proposals"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type SellerService struct {
	logger log.Logger
	db     *gorm.DB
	es     email.Sender
}

func NewSellerService(logger log.Logger, db *gorm.DB, es email.Sender) *SellerService {
	return &SellerService{logger: logger, db: db, es: es}
}

func (e SellerService) CreateProposal(ctx context.Context, userID string, createRequest proposals.CreateRequest) error {
	qID, _ := strconv.Atoi(createRequest.QuoteID)
	uID, _ := strconv.Atoi(userID)
	availableIntp := createRequest.AvailableDate
	availableInt := 0
	if availableIntp != nil {
		availableInt = *availableIntp
	}
	ad := time.Time{}
	ad = time.Unix(int64(availableInt), 0)

	freight, _ := strconv.ParseFloat(createRequest.Freight, 64)
	tax, _ := strconv.ParseFloat(createRequest.Tax, 64)
	fees, _ := strconv.ParseFloat(createRequest.Fees, 64)

	year, _ := strconv.Atoi(createRequest.Year)
	er := models.Proposal{
		ProposalNumber: strconv.Itoa(rangeIn(1000, 9999)),
		QuoteID:        qID,
		Make:           createRequest.Make,
		EqModel:        createRequest.EqModel,
		Year:           year,
		Serial:         createRequest.Serial,
		VIN:            createRequest.VIN,
		EqHours:        createRequest.EqHours,
		Freight:        freight,
		Tax: sql.NullFloat64{
			Float64: tax,
			Valid:   true,
		},
		Fees: sql.NullFloat64{
			Float64: fees,
			Valid:   true,
		},
		Condition: createRequest.Condition,
		SalePrice: createRequest.SalePrice,
		Description: sql.NullString{
			String: createRequest.Description,
			Valid:  true,
		},
		Comments: sql.NullString{
			String: createRequest.Comments,
			Valid:  true,
		},
		Specifications: sql.NullString{
			String: createRequest.Specifications,
			Valid:  true,
		},
		Videos: sql.NullString{
			String: createRequest.Videos,
			Valid:  true,
		},
		Pics: sql.NullString{
			String: createRequest.Pics,
			Valid:  true,
		},
		AvailableDate: sql.NullTime{
			Time:  ad,
			Valid: true,
		},
		WebUserID: uID,
		Status:    models.NewStatus,
		VendorCompanyName: createRequest.VendorCompanyName,
		VendorPhoneNumber: createRequest.VendorPhoneNumber,
		VendorEmail: createRequest.VendorEmail,
	}
	result := e.db.Create(&er)
	if result != nil && result.Error != nil {
		return result.Error
	}
	result = e.db.
		Find(&er, "id = ?", er.ID)
	if result.Error != nil {
		return result.Error
	}
	// err := e.notifyCreate(ctx, er)
	// if err != nil {
	// 	e.logger.Log("err", err.Error())
	// 	return err
	// }
	return nil
}

func (e SellerService) ShowProposal(ctx context.Context, proposalID string) (entities.Proposal, error) {
	var proposal models.Proposal
	result := e.db.
		Preload("WebUser").
		Preload("Quote.EquipmentRequest.Equipment").
		Find(&proposal, "id = ?", proposalID)
	if result.Error != nil {
		return entities.Proposal{}, result.Error
	}
	return proposal.ToCompleteProposalEntity(), nil
}

func (e SellerService) UpdateProposal(ctx context.Context, updateRequest proposals.UpdateRequest) error {
	var proposal models.Proposal
	result := e.db.
		Find(&proposal, "id = ?", updateRequest.ID)
	if result.Error != nil {
		return result.Error
	}
	amount := updateRequest.SalePrice
	// proposal.Zipcode = updateRequest.Zipcode
	proposal.SalePrice = amount
	// proposal.EquipmentRequest.SpecialRequest.Valid = true
	// proposal.EquipmentRequest.SpecialRequest.String = updateRequest.SpecialRequest
	// cp, _ := strconv.Atoi(updateRequest.ContractPreference)
	// proposal.EquipmentRequest.ContractPreference = cp
	// proposal.EquipmentRequest.RentFrom = sql.NullTime{
	// 	Time:  time.Time{},
	// 	Valid: false,
	// }
	// proposal.EquipmentRequest.RentTo = sql.NullTime{
	// 	Time:  time.Time{},
	// 	Valid: false,
	// }

	result = e.db.Save(&proposal)
	if result != nil && result.Error != nil {
		return result.Error
	}

	e.db.Model(&proposal).Where("id = ?", updateRequest.ID).Set("comments", updateRequest.Comments)
	// err := e.notifyChange(ctx, proposal)
	// if err != nil {
	// 	e.logger.Log("err", err.Error())
	// 	return err
	// }
	return nil
}

// func (e SellerService) notifyChange(ctx context.Context, proposal models.Proposal) error {
// 	var users []models.WebUser
// 	result := e.db.Find(&users, "profile_type = ?", "vendor_rentals")
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	emails := make([]string, 0)
// 	for _, u := range users {
// 		emails = append(emails, u.Username)
// 	}
// 	return e.es.SendChangeQuote(ctx, emails, proposal)
// }

// func (e SellerService) notifyCreate(ctx context.Context, quote models.Proposal) error {
// 	var users []models.WebUser
// 	result := e.db.Find(&users, "profile_type = ?", "vendor_rentals")
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	emails := make([]string, 0)
// 	for _, u := range users {
// 		emails = append(emails, u.Username)
// 	}
// 	e.logger.Log("emails", fmt.Sprintf("%v", emails))
// 	if len(emails) == 0 {
// 		return nil
// 	}
// 	return e.es.SendNewQuote(ctx, emails, []models.Quote{quote})
// }

func (e SellerService) DeleteProposal(ctx context.Context, proposalID string) error {
	var proposal models.Proposal
	result := e.db.Delete(&proposal, "id = ?", proposalID)
	if result != nil && result.Error != nil {
		return result.Error
	}
	return nil
}

func (u SellerService) Get(ctx context.Context, uqRequest proposals.SellerProposalRequest) ([]entities.Proposal, error) {
	var user models.WebUser
	result := u.db.Find(&user, "id = ?", uqRequest.UserID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}

	var proposals []models.Proposal
	result = u.db.
		Preload("WebUser").
		Preload("Quote.EquipmentRequest.Equipment").
		Find(&proposals, "web_user_id = ?", uqRequest.UserID).Where("proposals.status = ?", uqRequest.Status)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	eQ := make([]entities.Proposal, 0)
	for _, mq := range proposals {
		eQ = append(eQ, mq.ToCompleteProposalEntity())
	}
	return eQ, nil
}
func rangeIn(low, hi int) int {
	rand.Seed(time.Now().UnixNano())
	return low + rand.Intn(hi-low)
}
