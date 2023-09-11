package entities

type User struct {
	ID          string
	Password    string
	Email       string
	Phone       string
	Profile     Profile
	FullName    string
	CityName    string
	StateName   string
	CountryName string
	Address     string
	Zipcode     string
}

type UserInfo struct {
	ID                           int      `json:"id"`
	FullName                     string   `json:"full_name"`
	PhoneNumber                  string   `json:"phone_number"`
	Email                        string   `json:"email"`
	ProfilePicture               string   `json:"profile_picture"`
	AccountType                  int      `json:"account_type"`
	SublistNotification          bool     `json:"sublist_notifications"`
	EquipmentRequestNotification bool     `json:"equipment_notifications"`
	SupplyRequestNotification    bool     `json:"supplies_notifications"`
	EveryStateSelection          bool     `json:"every_state_selection"`
	NewUserPopUp                 bool     `json:"new_user_pop_up"`
	SupplyCategories             []uint   `json:"supply_categories"`
	EquipmentCategories          []uint   `json:"equipment_categories"`
	Cities                       []City   `json:"cities"`
	Counties                     []County `json:"counties"`
	States                       []State  `json:"states"`
	IsVerified                   bool     `json:"is_verified"`
	VerificationType             string   `json:"verification_type"`
}
