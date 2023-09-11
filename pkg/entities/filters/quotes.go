package filters

import "github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"

type ListQuotes struct {
	Keywords             []string
	Zipcode              string
	StateID              int
	CityID               int
	EquipmentCategories  []string
	SupplyCategories     []string
	ContractPreferences  []int
	RentTo               int64
	RentFrom             int64
	EquipmentSubcategory string
	UserID               string
	OnlySupply           bool
	OnlyEquipment        bool
	NotExpired           bool
	NotServed            bool
}

func MakeListQuotesFromRequest(listRequest requests.ListQuotesRequest) ListQuotes {
	return ListQuotes{
		Keywords:             listRequest.Keywords,
		Zipcode:              listRequest.Zipcode,
		StateID:              listRequest.StateID,
		CityID:               listRequest.CityID,
		EquipmentCategories:  listRequest.EquipmentCategories,
		SupplyCategories:     listRequest.SupplyCategories,
		ContractPreferences:  listRequest.ContractPreferences,
		RentTo:               listRequest.RentTo,
		RentFrom:             listRequest.RentFrom,
		EquipmentSubcategory: listRequest.EquipmentSubcategory,
		NotExpired:           true,
		NotServed:            true,
	}
}
