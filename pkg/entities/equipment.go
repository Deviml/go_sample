package entities

type Equipment struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Subcategory EquipmentSubcategory `json:"subcategory"`
}
