package entities

type Location struct {
	Zipcode string `json:"zipcode"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}
