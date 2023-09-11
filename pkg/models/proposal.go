package models

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"gorm.io/gorm"
)

const (
	NewStatus      = 0
	ApprovedStatus = 1
	RejectedStatus = 2
)

type Proposal struct {
	gorm.Model
	WebUser           WebUser
	WebUserID         int
	QuoteID           int
	Quote             *Quote
	Make              string
	EqModel           string
	Year              int
	Serial            string
	VIN               string
	EqHours           float32
	Condition         string
	SalePrice         float64
	AvailableDate     sql.NullTime
	Description       sql.NullString
	Comments          sql.NullString
	Specifications    sql.NullString // save comma separeted values
	Videos            sql.NullString // save comma separeted values
	Pics              sql.NullString // save comma separeted values
	Status            int
	Freight           float64
	Tax               sql.NullFloat64
	Fees              sql.NullFloat64
	ProposalNumber    string
	VendorCompanyName string
	VendorEmail       string
	VendorPhoneNumber string
}

func (s Proposal) TableName() string {
	return "proposals"
}

func (q Proposal) ToCompleteProposalEntity() entities.Proposal {
	// startDate := int64(0)
	// if q.EquipmentRequest.RentFrom.Valid {
	// 	startDate = q.EquipmentRequest.RentFrom.Time.Unix()
	// }
	// endDate := int64(0)
	// if q.EquipmentRequest.RentTo.Valid {
	// 	endDate = q.EquipmentRequest.RentTo.Time.Unix()
	// }
	return entities.Proposal{
		ProposalNumber: q.ProposalNumber,
		ID:             strconv.Itoa(int(q.ID)),
		Make:           q.Make,
		Model:          q.EqModel,
		Year:           strconv.Itoa(q.Year),
		Serial:         q.Serial,
		VIN:            q.VIN,
		EqHours:        fmt.Sprintf("%f", q.EqHours),
		Condition:      q.Condition,
		SalePrice:      fmt.Sprintf("%f", q.SalePrice),
		AvailableDate:  q.AvailableDate.Time.Unix(),
		Description:    q.Description.String,
		Comments:       q.Comments.String,
		Specifications: q.Specifications.String,
		Videos:         q.Videos.String,
		Pics:           q.Pics.String,
		Tax:            fmt.Sprintf("%f", q.Tax.Float64),
		Fees:           fmt.Sprintf("%f", q.Fees.Float64),
		Status:         q.GetProposalStatus(),
		Freight:        fmt.Sprintf("%f", q.Freight),
		EquipmentName:  q.Quote.EquipmentRequest.Equipment.Name,
		VendorCompanyName: q.VendorCompanyName,
		VendorEmail:       q.VendorEmail,
		VendorPhoneNumber: q.VendorPhoneNumber,
		// EquipmentCategory:  q.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory.ToEntity(),
		// Equipment:          q.EquipmentRequest.Equipment.ToEntity(),
		// ContractPreference: q.EquipmentRequest.ContractPreference,
		// StartDateFormatted: q.EquipmentRequest.RentFrom.Time,
		// StartDate:          startDate,
		// EndDateFormatted:   q.EquipmentRequest.RentTo.Time,
		// EndDate:            endDate,
		// Amount:             q.EquipmentRequest.Amount,
		// Subcontractor: entities.SubcontractorSummary{
		// 	Location: entities.Location{
		// 		Zipcode: q.Zipcode,
		// 		City:    q.City.Name,
		// 		State:   q.City.State.Name,
		// 		Country: "USA",
		// 	},
		// 	Company: entities.Company{
		// 		ID:   "1",
		// 		Name: "Request Test Company",
		// 	},
		// 	Email:    q.WebUser.Username,
		// 	Phone:    q.WebUser.Phone,
		// 	FullName: q.WebUser.FullName,
		// },
		// SpecialRequest: q.EquipmentRequest.SpecialRequest.String,
		// PurchasedAt:    q.PurchasedAt.Unix(),
	}
}
func (p Proposal) GetProposalStatus() string {
	switch p.Status {
	case ApprovedStatus:
		return "Approved"
	case RejectedStatus:
		return "Rejected"
	default:
		return "Pending"
	}
}
