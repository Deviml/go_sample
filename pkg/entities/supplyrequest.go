package entities

type SupplyRequest struct {
	Name   string `json:"name"`
	Supply Supply `json:"supply"`
	Amount int    `json:"amount"`
}

type CompleteSupplyRequest struct {
	ID             string               `json:"id"`
	Name           string               `json:"name"`
	SupplyCategory SupplyCategory       `json:"supply_category"`
	Location       Location             `json:"location"`
	Supply         Supply               `json:"supply"`
	Subcontractor  SubcontractorSummary `json:"subcontractor"`
	Amount         int                  `json:"amount"`
	SpecialRequest string               `json:"special_request"`
	PurchasedAt    int64                `json:"purchased_at"`
}
