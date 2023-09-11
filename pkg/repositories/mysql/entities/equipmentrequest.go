package entities

import (
	"database/sql"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
)

type EquipmentRequest struct {
	ID                 string         `json:"id"`
	Name               sql.NullString `json:"name"`
	Equipment          Equipment      `json:"equipment"`
	ContractPreference sql.NullInt32  `json:"contract_preference"`
	StartDateFormatted sql.NullTime   `json:"start_date_formatted"`
	EndDateFormatted   sql.NullTime   `json:"end_date_formatted"`
	ExpirationDate     sql.NullTime   `json:"expiration_date"`
	Amount             sql.NullInt32  `json:"amount"`
	SpecialRequest     sql.NullString `json:"special_request"`
}

func (e Equipment) ToEntity() entities.Equipment {
	return entities.Equipment{
		ID:          e.ID,
		Name:        e.Name.String,
		Subcategory: e.EquipmentSubcategory.ToEntity(),
	}
}
