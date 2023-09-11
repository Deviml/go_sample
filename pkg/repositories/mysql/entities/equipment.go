package entities

import "database/sql"

type Equipment struct {
	ID                   string               `json:"id"`
	Name                 sql.NullString       `json:"name"`
	EquipmentSubcategory EquipmentSubcategory `json:"equipment_subcategory"`
}
