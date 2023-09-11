package entities

type ProfileType string

const (
	VendorRentalProfileType      ProfileType = "1"
	SubcontractorProfileType     ProfileType = "2"
	GeneralContractorProfileType ProfileType = "3"
)

type Profile struct {
	ID   string
	Type ProfileType
}

type VendorRental struct {
	Profile
	StateID string
	CityID  string
	Company Company
}

type Subcontractor struct {
	Profile
	StateID string
	CityID  string
	Company Company
}

type CompleteSub struct {
	Location Location `json:"location"`
	Company  Company  `json:"company"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	FullName string   `json:"full_name"`
}

type GeneralContractor struct {
	Profile
	Name             string  `json:"name"`
	Company          Company `json:"company"`
	StateID          string
	CityID           string
	VerificationCode string
}

type GeneralContractorSummary struct {
	Name    string  `json:"name"`
	Company Company `json:"company"`
}

type SubcontractorSummary struct {
	Location Location `json:"location"`
	Company  Company  `json:"company"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	FullName string   `json:"full_name"`
}
