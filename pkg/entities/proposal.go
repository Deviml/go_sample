package entities

type Proposal struct {
	ID                string               `json:"id"`
	Make              string               `json:"make"`
	Model             string               `json:"eq_model"`
	Year              string               `json:"year"`
	Serial            string               `json:"serial"`
	VIN               string               `json:"vin"`
	EqHours           string               `json:"eq_hours"`
	Condition         string               `json:"condition"`
	SalePrice         string               `json:"sale_price"`
	AvailableDate     int64                `json:"available_date"`
	Description       string               `json:"description"`
	Comments          string               `json:"comments"`
	Specifications    string               `json:"specifications"`
	Videos            string               `json:"videos"`
	Pics              string               `json:"pics"`
	Status            string               `json:"status"`
	EquipmentName     string               `json:"equipment_name"`
	Subcontractor     SubcontractorSummary `json:"subcontractor"`
	Freight           string               `json:"freight"`
	Tax               string               `json:"tax"`
	Fees              string               `json:"fees"`
	ProposalNumber    string               `json:"proposal_number"`
	VendorCompanyName string               `json:"vendor_company_name"`
	VendorEmail       string               `json:"vendor_email"`
	VendorPhoneNumber string               `json:"vendor_phone_number"`
}
