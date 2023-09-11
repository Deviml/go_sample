package entities

type Supply struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	SupplyCategory SupplyCategory `json:"supply_category"`
}
