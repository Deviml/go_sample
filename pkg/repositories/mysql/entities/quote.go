package entities

import (
	"database/sql"
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
)

const (
	quoteTypeEquipment  = 1
	quoteTypeSupply     = 2
	requestTypeRent     = 1
	requestTypePurchase = 2
)

type Quote struct {
	ID               string            `json:"id"`
	Type             sql.NullInt32     `json:"type"`
	RequestType      sql.NullInt32     `json:"request_type"`
	EquipmentRequest EquipmentRequest  `json:"equipment_request"`
	SupplyRequest    SupplyRequest     `json:"supply_request"`
	Description      sql.NullString    `json:"description"`
	Location         entities.Location `json:"location"`
	ProductName      string            `json:"product_name"`
	WebUser          models.WebUser    `json:"-"`
	PurchasedAt      time.Time         `json:"-"`
	Status           int               `json:"status"`
}

func (q Quote) ToCompleteSupplyRequestEntity() entities.CompleteSupplyRequest {
	return entities.CompleteSupplyRequest{
		ID:             q.ID,
		Name:           q.SupplyRequest.Supply.Name.String,
		SupplyCategory: q.SupplyRequest.Supply.SupplyCategory.ToEntity(),
		Location:       q.Location,
		Supply:         q.SupplyRequest.Supply.ToEntity(),
		Subcontractor: entities.SubcontractorSummary{
			Location: entities.Location{
				Zipcode: q.Location.Zipcode,
				City:    q.Location.City,
				State:   q.Location.State,
				Country: "USA",
			},
			Company: entities.Company{
				ID:   "1",
				Name: "Request Test Company",
			},
			Email:    q.WebUser.Username,
			Phone:    q.WebUser.Phone,
			FullName: q.WebUser.FullName,
		},
		Amount:         int(q.SupplyRequest.Amount.Int32),
		SpecialRequest: q.SupplyRequest.SpecialRequest.String,
		PurchasedAt:    q.PurchasedAt.Unix(),
	}
}

func (q Quote) ToCompleteEquipmentRequestEntity() entities.CompleteEquipmentRequest {
	startDate := int64(0)
	if q.EquipmentRequest.StartDateFormatted.Valid {
		startDate = q.EquipmentRequest.StartDateFormatted.Time.Unix()
	}
	endDate := int64(0)
	if q.EquipmentRequest.EndDateFormatted.Valid {
		endDate = q.EquipmentRequest.EndDateFormatted.Time.Unix()
	}
	return entities.CompleteEquipmentRequest{
		ID:                 q.ID,
		Name:               q.EquipmentRequest.Equipment.Name.String,
		EquipmentCategory:  q.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory.ToEntity(),
		Equipment:          q.EquipmentRequest.Equipment.ToEntity(),
		ContractPreference: int(q.EquipmentRequest.ContractPreference.Int32),
		StartDateFormatted: q.EquipmentRequest.StartDateFormatted.Time,
		StartDate:          startDate,
		EndDateFormatted:   q.EquipmentRequest.EndDateFormatted.Time,
		EndDate:            endDate,
		Amount:             int(q.EquipmentRequest.Amount.Int32),
		Subcontractor: entities.SubcontractorSummary{
			Location: entities.Location{
				Zipcode: q.Location.Zipcode,
				City:    q.Location.City,
				State:   q.Location.State,
				Country: "USA",
			},
			Company: entities.Company{
				ID:   "1",
				Name: "Request Test Company",
			},
			Email:    q.WebUser.Username,
			Phone:    q.WebUser.Phone,
			FullName: q.WebUser.FullName,
		},
		SpecialRequest: q.EquipmentRequest.SpecialRequest.String,
		PurchasedAt:    q.PurchasedAt.Unix(),
		Status:         q.GetQuoteStatus(),
	}
}

func (q Quote) ToDomainQuote() entities.Quote {
	quoteType := quoteTypeSupply
	if q.EquipmentRequest.Name.Valid {
		quoteType = quoteTypeEquipment
	}

	requestType := requestTypePurchase

	if q.EquipmentRequest.StartDateFormatted.Valid && q.EquipmentRequest.EndDateFormatted.Valid {
		requestType = requestTypeRent
	}

	var equipment *entities.EquipmentRequest

	if q.EquipmentRequest.Equipment.Name.Valid {
		equipment = &entities.EquipmentRequest{
			Name: q.EquipmentRequest.Name.String,
			Equipment: entities.Equipment{
				ID:   q.EquipmentRequest.Equipment.ID,
				Name: q.EquipmentRequest.Equipment.Name.String,
				Subcategory: entities.EquipmentSubcategory{
					ID:   q.EquipmentRequest.Equipment.EquipmentSubcategory.ID,
					Name: q.EquipmentRequest.Equipment.EquipmentSubcategory.Name.String,
					EquipmentCategory: entities.EquipmentCategory{
						ID:   q.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory.ID,
						Name: q.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory.Name.String,
					},
				},
			},
			ContractPreference: int(q.EquipmentRequest.ContractPreference.Int32),
			StartDateFormatted: q.EquipmentRequest.StartDateFormatted.Time,
			StartDate:          q.EquipmentRequest.StartDateFormatted.Time.Unix(),
			EndDateFormatted:   q.EquipmentRequest.EndDateFormatted.Time,
			EndDate:            q.EquipmentRequest.EndDateFormatted.Time.Unix(),
			ExpirationDate:     q.EquipmentRequest.ExpirationDate.Time.Unix(),
			Amount:             int(q.EquipmentRequest.Amount.Int32),
			SpecialRequest:     q.EquipmentRequest.SpecialRequest.String,
		}
	}

	var supply *entities.SupplyRequest

	if q.SupplyRequest.Supply.Name.Valid {
		supply = &entities.SupplyRequest{
			Name: q.SupplyRequest.Name.String,
			Supply: entities.Supply{
				ID:   q.SupplyRequest.Supply.ID,
				Name: q.SupplyRequest.Supply.Name.String,
				SupplyCategory: entities.SupplyCategory{
					ID:   q.SupplyRequest.Supply.SupplyCategory.ID,
					Name: q.SupplyRequest.Supply.SupplyCategory.Name.String,
				},
			},
			Amount: int(q.SupplyRequest.Amount.Int32),
		}
	}

	description := q.EquipmentRequest.Name.String
	if description == "" {
		description = q.SupplyRequest.Name.String
	}

	return entities.Quote{
		ID:               q.ID,
		Type:             quoteType,
		RequestType:      requestType,
		EquipmentRequest: equipment,
		SupplyRequest:    supply,
		Description:      description,
		Location:         q.Location,
		ExpiryDate:       q.EquipmentRequest.ExpirationDate.Time.Unix(),
		Status:           q.GetQuoteStatus(),
	}
}
func (q Quote) GetQuoteStatus() string {
	switch q.Status {
	case models.PurchasedStatus:
		return "Purchased"
	case models.ServedStatus:
		return "Served"
	case models.NewStatus:
		return "New"
	default:
		return "Expired"
	}
}
