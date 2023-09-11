package entities

type Company struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CompanyNew struct {
	ID string `json:"id"`
}
