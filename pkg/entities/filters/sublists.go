package filters

type ListSublists struct {
	StateID  int
	Keywords []string
	Zipcode  string
	CityID   int
	From     string
	To       string
	UserID   string
}
