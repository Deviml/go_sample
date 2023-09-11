package entities

type EquipmentSubcategory struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	EquipmentCategory EquipmentCategory `json:"equipment_category"`
}
