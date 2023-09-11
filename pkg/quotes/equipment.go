package quotes

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/quotes"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type EquipmentService struct {
	logger log.Logger
	db     *gorm.DB
	es     email.Sender
}

func NewEquipmentService(logger log.Logger, db *gorm.DB, es email.Sender) *EquipmentService {
	return &EquipmentService{logger: logger, db: db, es: es}
}

func (e EquipmentService) CreateEquipment(ctx context.Context, userID string, createRequest quotes.CreateRequest) error {
	eID, _ := strconv.Atoi(createRequest.EquipmentID)
	uID, _ := strconv.Atoi(userID)
	locationID := 1
	endIntp := createRequest.EndDate
	endInt := 0
	if endIntp != nil {
		endInt = *endIntp
	}
	et := time.Time{}
	et = time.Unix(int64(endInt), 0)

	amount, _ := strconv.Atoi(createRequest.Amount)
	cID, _ := strconv.Atoi(createRequest.CityID)
	_, _ = strconv.Atoi(createRequest.CountyID)
	er := models.Quote{
		EquipmentRequest: &models.EquipmentRequest{
			SpecialRequest: sql.NullString{
				String: createRequest.SpecialRequest,
				Valid:  true,
			},
			EquipmentID:        uint(eID),
			Amount:             amount,
			LocationID:         locationID,
			ContractPreference: createRequest.ContractPreference,
			RentFrom: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			RentTo: sql.NullTime{
				Time:  et,
				Valid: true,
			},
			ExpirationDate: sql.NullTime{
				Time:  time.Now().AddDate(0, 0, 14),
				Valid: true,
			},
		},
		WebUserID:      uID,
		CityID:         cID,
		Zipcode:        createRequest.Zipcode,
		Address:        createRequest.Address,
		CountyID:       nil,
		Status:         models.NewStatus,
		ExpirationDate: time.Now().AddDate(0, 0, 14),
	}
	result := e.db.Create(&er)
	if result != nil && result.Error != nil {
		return result.Error
	}
	result = e.db.
		Preload("EquipmentRequest.Location").
		Preload("City.State").
		Preload("EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory").
		Find(&er, "id = ?", er.ID)
	if result.Error != nil {
		return result.Error
	}
	err := e.notifyCreate(ctx, er)
	if err != nil {
		e.logger.Log("err", err.Error())
		return err
	}
	return nil
}
func (e EquipmentService) ShowEquipment(ctx context.Context, quoteID string) (entities.CompleteEquipmentRequest, error) {
	var quote models.Quote
	result := e.db.
		Preload("EquipmentRequest.Location").
		Preload("City.State").
		Preload("EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory").
		Find(&quote, "id = ?", quoteID)
	if result.Error != nil {
		return entities.CompleteEquipmentRequest{}, result.Error
	}
	return quote.ToCompleteEquipmentRequestEntity(), nil
}

func (e EquipmentService) UpdateEquipment(ctx context.Context, updateRequest quotes.UpdateRequest) error {
	var quote models.Quote

	result := e.db.Model(&quote).Where("id = ?", updateRequest.ID).Update("zipcode", updateRequest.Zipcode)
	result = e.db.
		Preload("EquipmentRequest.Location").
		Preload("City.State").
		Preload("EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory").
		Find(&quote, "id = ?", updateRequest.ID)
	if result.Error != nil {
		return result.Error
	}
	amount := updateRequest.Amount
	quote.EquipmentRequest.Amount = amount
	quote.EquipmentRequest.SpecialRequest.Valid = true
	quote.EquipmentRequest.SpecialRequest.String = updateRequest.SpecialRequest
	cp, _ := strconv.Atoi(updateRequest.ContractPreference)
	quote.EquipmentRequest.ContractPreference = cp

	result = e.db.Save(&quote.EquipmentRequest)
	if result != nil && result.Error != nil {
		return result.Error
	}
	return nil
}

func (e EquipmentService) notifyCreate(ctx context.Context, quote models.Quote) error {
	var users []models.WebUser
	result := e.db.Preload("Cities").Preload("Counties").Preload("States").Preload("EquipmentCategories").Find(&users, "profile_type = ?", "vendor_rentals")
	if result.Error != nil {
		return result.Error
	}
	emails := make([]string, 0)
	for _, u := range users {
		if !u.EquipmentRequestNotification {
			e.logger.Log("skip_user_equipment", u.ID)
			continue
		}
		e.logger.Log("user_equipment_categories", fmt.Sprintf("%v", u.EquipmentCategories))
		e.logger.Log("quote_category", quote.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory.ID)
		if !u.HaveEquipmentCategory(quote.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory.ID) {
			e.logger.Log("skip_user_equipment_category", u.ID)
			continue
		}

		if len(u.States) > 0 {
			if len(u.Cities) > 0 {
				if !u.HaveCity(uint(quote.City.ID)) {
					e.logger.Log("skip_user", u.ID)
					continue
				}
			}
			if !u.HaveState(quote.City.StateID) {
				e.logger.Log("skip_user", u.ID)
				continue
			}
		} else {
			e.logger.Log("skip_user_no_preference", u.ID)
			continue
		}
		emails = append(emails, u.Username)
	}
	e.logger.Log("emails", fmt.Sprintf("%v", emails))
	if len(emails) == 0 {
		return nil
	}

	Name := quote.EquipmentRequest.Equipment.Name
	Category := quote.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory.Name
	Amount := quote.EquipmentRequest.Amount
	City := quote.City.Name
	Zipcode := quote.Zipcode
	SpecialRequest := quote.EquipmentRequest.SpecialRequest.String
	State := quote.City.State.Name
	ID := quote.ID
	FrontURL := "www.equiphunter.com"

	return e.es.SendNewQuote(ctx, emails, Name, Category, Amount, City, Zipcode, SpecialRequest, State, FrontURL, ID)
}

func (e EquipmentService) DeleteEquipment(ctx context.Context, quoteID string) error {
	var quote models.Quote
	result := e.db.Delete(&quote, "id = ?", quoteID)
	if result != nil && result.Error != nil {
		return result.Error
	}
	return nil
}
