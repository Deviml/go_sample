package entities

import "time"

type EquipmentRequest struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Equipment          Equipment `json:"equipment"`
	ContractPreference int       `json:"contract_preference"`
	StartDateFormatted time.Time `json:"start_date_formatted"`
	StartDate          int64     `json:"start_date"`
	EndDateFormatted   time.Time `json:"end_date_formatted"`
	EndDate            int64     `json:"end_date"`
	ExpirationDate     int64     `json:"expiration_date"`
	Amount             int       `json:"amount"`
	SpecialRequest     string    `json:"special_request"`
}

type CompleteEquipmentRequest struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	EquipmentCategory  EquipmentCategory    `json:"equipment_category"`
	Equipment          Equipment            `json:"equipment"`
	ContractPreference int                  `json:"contract_preference"`
	StartDateFormatted time.Time            `json:"start_date_formatted"`
	StartDate          int64                `json:"start_date"`
	EndDateFormatted   time.Time            `json:"end_date_formatted"`
	EndDate            int64                `json:"end_date"`
	Amount             int                  `json:"amount"`
	Subcontractor      SubcontractorSummary `json:"subcontractor"`
	SpecialRequest     string               `json:"special_request"`
	Zipcode            string               `json:"zipcode"`
	Address            string               `json:"address"`
	PurchasedAt        int64                `json:"purchased_at"`
	Status             string               `json:"status"`
}
