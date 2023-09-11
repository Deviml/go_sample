package quotes

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/quotes"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type SupplyService struct {
	logger log.Logger
	db     *gorm.DB
	es     email.Sender
}

func NewSupplyService(logger log.Logger, db *gorm.DB, es email.Sender) *SupplyService {
	return &SupplyService{logger: logger, db: db, es: es}
}

func (s SupplyService) CreateSupply(ctx context.Context, userID string, createRequest quotes.SupplyCreateRequest) error {
	sID, _ := strconv.Atoi(createRequest.SupplyID)
	uID, _ := strconv.Atoi(userID)
	locationID := 1
	amount, _ := strconv.Atoi(createRequest.Amount)
	cID, _ := strconv.Atoi(createRequest.CityID)
	_, _ = strconv.Atoi(createRequest.CountyID)
	quote := models.Quote{
		WebUserID: uID,
		SupplyRequest: &models.SupplyRequest{
			SpecialRequest: sql.NullString{
				String: createRequest.SpecialRequest,
				Valid:  true,
			},
			SupplyID:   uint(sID),
			Amount:     amount,
			LocationID: locationID,
		},
		CityID:   cID,
		CountyID: nil,
		Zipcode:  createRequest.Zipcode,
		Status:   models.NewStatus,
	}
	result := s.db.Create(&quote)
	if result != nil && result.Error != nil {
		return result.Error
	}
	result = s.db.
		Preload("SupplyRequest.Location").
		Preload("City.State").
		Preload("SupplyRequest.Supply.SupplyCategory").Find(&quote, "id = ?", quote.ID)
	if result.Error != nil {
		return result.Error
	}
	err := s.notifyCreate(ctx, quote)
	if err != nil {
		s.logger.Log("err", err.Error())
		return err
	}
	return nil
}

func (s SupplyService) notifyCreate(ctx context.Context, quote models.Quote) error {
	var users []models.WebUser
	result := s.db.Preload("Cities").Preload("Counties").Preload("States").Preload("SupplyCategories").Find(&users, "profile_type = ?", "vendor_rentals")
	if result.Error != nil {
		return result.Error
	}
	emails := make([]string, 0)
	for _, u := range users {
		if !u.SupplyRequestNotification {
			continue
		}
		if !u.HaveSupplyCategory(quote.SupplyRequest.Supply.SupplyCategory.ID) {
			continue
		}

		if len(u.States) > 0 {
			if len(u.Cities) > 0 {
				if !u.HaveCity(uint(quote.City.ID)) {
					s.logger.Log("skip_user", u.ID)
					continue
				}
			}
			if !u.HaveState(quote.City.StateID) {
				s.logger.Log("skip_user", u.ID)
				continue
			}
		} else {
			s.logger.Log("skip_user_no_preference", u.ID)
			continue
		}
		emails = append(emails, u.Username)
	}
	s.logger.Log("email", fmt.Sprintf("%v", emails))
	if len(emails) == 0 {
		return nil
	}
	// CHANGED FOR NOW, DURING DEVELOPMENT OF BULK EMAIL NOTIFICATIONS, WILL CHANGE BACK WHEN NEEDED BY SYSTEM AS WE ARE NOT USING SUPPLU AND SUBLISTS FOR NOW
	return nil
}

func (s SupplyService) ShowSupply(ctx context.Context, quoteID string) (entities.CompleteSupplyRequest, error) {
	var quote models.Quote
	result := s.db.
		Preload("SupplyRequest.Location").
		Preload("City.State").
		Preload("SupplyRequest.Supply.SupplyCategory").Find(&quote, "id = ?", quoteID)
	if result.Error != nil {
		return entities.CompleteSupplyRequest{}, result.Error
	}
	return quote.ToCompleteSupplyRequestEntity(), nil
}

func (s SupplyService) UpdateSupply(ctx context.Context, updateRequest quotes.SupplyUpdateRequest) error {
	var quote models.Quote
	result := s.db.
		Preload("SupplyRequest.Location").
		Preload("City.State").
		Preload("SupplyRequest.Supply.SupplyCategory").Find(&quote, "id = ?", updateRequest.ID)
	if result.Error != nil {
		return result.Error
	}
	amount := updateRequest.Amount
	quote.SupplyRequest.Amount = amount
	quote.SupplyRequest.SpecialRequest.Valid = true
	quote.SupplyRequest.SpecialRequest.String = updateRequest.SpecialRequest
	result = s.db.Save(&quote.SupplyRequest)
	if result != nil && result.Error != nil {
		return result.Error
	}

	result = s.db.Model(&quote).Where("id = ?", updateRequest.ID).Update("zipcode", updateRequest.Zipcode)
	if result != nil && result.Error != nil {
		return result.Error
	}
	return nil
}

func (s SupplyService) DeleteSupply(ctx context.Context, quoteID string) error {
	var quote models.Quote
	result := s.db.Delete(&quote, "id = ?", quoteID)
	if result != nil && result.Error != nil {
		return result.Error
	}
	return nil
}
