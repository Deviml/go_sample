package entities

import (
	"database/sql"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
)

type EquipmentCategory struct {
	ID   string         `json:"id"`
	Name sql.NullString `json:"name"`
}

func (e EquipmentCategory) ToEntity() entities.EquipmentCategory {
	return entities.EquipmentCategory{
		ID:   e.ID,
		Name: e.Name.String,
	}
}
