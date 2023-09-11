package entities

import (
	"database/sql"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
)

type EquipmentSubcategory struct {
	ID                string            `json:"id"`
	Name              sql.NullString    `json:"name"`
	EquipmentCategory EquipmentCategory `json:"equipment_category"`
}

func (e EquipmentSubcategory) ToEntity() entities.EquipmentSubcategory {
	return entities.EquipmentSubcategory{
		ID:                e.ID,
		Name:              e.Name.String,
		EquipmentCategory: e.EquipmentCategory.ToEntity(),
	}
}
