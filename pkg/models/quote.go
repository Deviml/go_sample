package models

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"gorm.io/gorm"
)

const (
	EquipmentRequestPurchase = "1"
	EquipmentRequestRent     = "2"
	PurchasedStatus          = 1
	ServedStatus             = 2
	ExpiredStatus            = 3
)

type SupplyRequest struct {
	gorm.Model
	SupplyID       uint `gorm:"type:int(11);"`
	Supply         Supply
	Amount         int
	LocationID     int `gorm:"type:int(10);"`
	Location       Location
	SpecialRequest sql.NullString
}

func (s SupplyRequest) TableName() string {
	return "SupplyRequests"
}

type EquipmentRequest struct {
	gorm.Model
	EquipmentID        uint
	Equipment          Equipment
	Amount             int
	LocationID         int
	Location           Location
	ContractPreference int
	RentFrom           sql.NullTime
	RentTo             sql.NullTime
	ExpirationDate     sql.NullTime
	SpecialRequest     sql.NullString
}

func (e EquipmentRequest) TableName() string {
	return "EquipmentRequests"
}

type Quote struct {
	gorm.Model
	SupplyRequestID    *uint `gorm:"type:int(11)"`
	SupplyRequest      *SupplyRequest
	EquipmentRequestID *uint `gorm:"type:int(11)"`
	EquipmentRequest   *EquipmentRequest
	WebUser            WebUser
	WebUserID          int
	VendorRentals      []VendorRental `gorm:"many2many:vendor_rental_quotes;"`
	CityID             int
	City               *City
	Zipcode            string
	Address            string
	CountyID           *int
	County             *County
	PurchasedAt        time.Time `gorm:"-"`
	Status             int
	ExpirationDate     time.Time
}

func (q Quote) TableName() string {
	return "Quotes"
}

func (q Quote) ToCompleteSupplyRequestEntity() entities.CompleteSupplyRequest {
	return entities.CompleteSupplyRequest{
		ID:             strconv.Itoa(int(q.ID)),
		Name:           q.SupplyRequest.Supply.Name,
		SupplyCategory: q.SupplyRequest.Supply.SupplyCategory.ToEntity(),
		Location:       q.SupplyRequest.Location.ToEntity(),
		Supply:         q.SupplyRequest.Supply.ToEntity(),
		Subcontractor: entities.SubcontractorSummary{
			Location: entities.Location{
				Zipcode: q.Zipcode,
				City:    q.City.Name,
				State:   q.City.State.Name,
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
		Amount:         q.SupplyRequest.Amount,
		SpecialRequest: q.SupplyRequest.SpecialRequest.String,
		PurchasedAt:    q.PurchasedAt.Unix(),
	}
}

func (q Quote) ToCompleteEquipmentRequestEntity() entities.CompleteEquipmentRequest {
	startDate := int64(0)
	if q.EquipmentRequest.RentFrom.Valid {
		startDate = q.EquipmentRequest.RentFrom.Time.Unix()
	}
	endDate := int64(0)
	if q.EquipmentRequest.RentTo.Valid {
		endDate = q.EquipmentRequest.RentTo.Time.Unix()
	}
	return entities.CompleteEquipmentRequest{
		ID:                 strconv.Itoa(int(q.ID)),
		Name:               q.EquipmentRequest.Equipment.Name,
		EquipmentCategory:  q.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory.ToEntity(),
		Equipment:          q.EquipmentRequest.Equipment.ToEntity(),
		ContractPreference: q.EquipmentRequest.ContractPreference,
		StartDateFormatted: q.EquipmentRequest.RentFrom.Time,
		StartDate:          startDate,
		EndDateFormatted:   q.EquipmentRequest.RentTo.Time,
		EndDate:            endDate,
		Amount:             q.EquipmentRequest.Amount,
		Subcontractor: entities.SubcontractorSummary{
			Location: entities.Location{
				Zipcode: q.Zipcode,
				City:    q.City.Name,
				State:   q.City.State.Name,
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
		Zipcode:        q.Zipcode,
		Address:        q.Address,
		PurchasedAt:    q.PurchasedAt.Unix(),
		Status:         q.GetQuoteStatus(),
	}
}

func (q Quote) ToQuiteEntity() entities.CompleteQuote {
	if q.EquipmentRequest != nil {
		ec := q.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory.ToEntity()
		return entities.CompleteQuote{
			ID:                 strconv.Itoa(int(q.ID)),
			Name:               q.EquipmentRequest.Equipment.Name,
			EquipmentCategory:  &ec,
			ContractPreference: q.EquipmentRequest.ContractPreference,
			StartDate:          q.EquipmentRequest.RentFrom.Time.Unix(),
			EndDate:            q.EquipmentRequest.RentTo.Time.Unix(),
			ExpiryDate:         q.EquipmentRequest.ExpirationDate.Time.Unix(),
			Amount:             q.EquipmentRequest.Amount,
			Subcontractor: entities.CompleteSub{
				Location: entities.Location{
					Zipcode: q.Zipcode,
					City:    q.City.Name,
					State:   q.City.State.Name,
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
		}
	}
	sc := q.SupplyRequest.Supply.SupplyCategory.ToEntity()
	return entities.CompleteQuote{
		ID:             strconv.Itoa(int(q.ID)),
		Name:           q.SupplyRequest.Supply.Name,
		SupplyCategory: &sc,
		Subcontractor: entities.CompleteSub{
			Location: entities.Location{
				Zipcode: q.Zipcode,
				City:    q.City.Name,
				State:   q.City.State.Name,
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
		Amount: q.SupplyRequest.Amount,
	}
}

func (q Quote) ToSingleQuoteEntity() entities.CompleteSingleQuote {
	return entities.CompleteSingleQuote{
		ID:                   strconv.Itoa(int(q.ID)),
		EquipmentName:        q.EquipmentRequest.Equipment.Name,
		EquipmentSubCategory: q.EquipmentRequest.Equipment.EquipmentSubcategory.Name,
		ContractPreference:   q.EquipmentRequest.ContractPreference,
		EndDate:              q.EquipmentRequest.RentTo.Time.Unix(),
		ExpiryDate:           q.EquipmentRequest.ExpirationDate.Time.Unix(),
		Amount:               q.EquipmentRequest.Amount,
		SpecialRequest:       q.EquipmentRequest.SpecialRequest.String,
		Subcontractor: entities.CompleteSub{
			Location: entities.Location{
				Zipcode: q.Zipcode,
				City:    q.City.Name,
				State:   q.City.State.Name,
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
	}
}

func (q Quote) GetQuoteStatus() string {
	switch q.Status {
	case PurchasedStatus:
		return "Purchased"
	case ServedStatus:
		return "Served"
	case NewStatus:
		return "New"
	default:
		return "Expired"
	}
}
