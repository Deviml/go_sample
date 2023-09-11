package entities

import (
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
)

type Proposal struct {
	ID             string         `json:"id"`
	Make           string         `json:"make"`
	Model          string         `json:"eq_model"`
	Year           string         `json:"year"`
	Serial         string         `json:"serial"`
	Vin            string         `json:"vin"`
	EqHours        string         `json:"eq_hours"`
	Condition      string         `json:"condition"`
	SalePrice      string         `json:"sale_price"`
	AvailableDate  time.Time      `json:"available_date"`
	Description    string         `json:"description"`
	Comments       string         `json:"comments"`
	Specifications string         `json:"specifications"`
	Videos         string         `json:"videos"`
	Pics           string         `json:"pics"`
	Status         string         `json:"status"`
	EquipmentName  string         `json:"equipment_name"`
	WebUser        models.WebUser `json:"-"`
	Freight        string         `json:"freight"`
	Tax            string         `json:"tax"`
	Fees           string         `json:"fees"`
	ProposalNumber string         `json:"proposal_number"`
}

func (q Proposal) ToListProposals() entities.Proposal {

	return entities.Proposal{
		ProposalNumber: q.ProposalNumber,
		ID:             q.ID,
		Make:           q.Make,
		Model:          q.Model,
		Year:           q.Year,
		Serial:         q.Serial,
		VIN:            q.Vin,
		EqHours:        q.EqHours,
		Condition:      q.Condition,
		SalePrice:      q.SalePrice,
		AvailableDate:  q.AvailableDate.Unix(),
		Description:    q.Description,
		Comments:       q.Comments,
		Specifications: q.Specifications,
		Videos:         q.Videos,
		Pics:           q.Pics,
		Status:         q.Status,
		EquipmentName:  q.EquipmentName,
		Freight:        q.Freight,
		Tax:            q.Tax,
		Fees:           q.Fees,
	}
}
