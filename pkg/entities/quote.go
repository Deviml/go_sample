package entities

type Quote struct {
	ID               string            `json:"id"`
	Type             int               `json:"type"`
	RequestType      int               `json:"request_type"`
	EquipmentRequest *EquipmentRequest `json:"equipment_request"`
	SupplyRequest    *SupplyRequest    `json:"supply_request"`
	Description      string            `json:"description"`
	Location         Location          `json:"location"`
	Amount           string            `json:"amount"`
	Status           string            `json:"status"`
	ExpiryStatus     string            `json:"expiry_status"`
	ExpiryDate       int64             `json:"expiry_date"`
}

type CompleteQuote struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	EquipmentCategory  *EquipmentCategory `json:"equipment_category"`
	ContractPreference int                `json:"contract_preference"`
	StartDate          int64              `json:"start_date"`
	EndDate            int64              `json:"end_date"`
	Subcontractor      CompleteSub        `json:"subcontractor"`
	SupplyCategory     *SupplyCategory    `json:"supply_category"`
	Amount             int                `json:"amount"`
	SpecialRequest     string             `json:"special_request"`
	Status             string             `json:"status"`
	ExpiryDate         int64              `json:"expiry_date"`
}

type CompleteSingleQuote struct {
	ID                   string      `json:"id"`
	EquipmentSubCategory string      `json:"equipment_sub_category"`
	EquipmentName        string      `json:"equipment_name"`
	ContractPreference   int         `json:"contract_preference"`
	Amount               int         `json:"amount"`
	EndDate              int64       `json:"end_date"`
	ExpiryDate           int64       `json:"expiry_date"`
	SpecialRequest       string      `json:"special_request"`
	Subcontractor        CompleteSub `json:"subcontractor"`
}
