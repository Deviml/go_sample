package requests

import (
	"context"
	"errors"
	"net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/helpers"
	"github.com/gorilla/mux"
)

type ListQuotesRequest struct {
	Keywords             []string `json:"keywords"`
	Zipcode              string   `json:"zipcode"`
	StateID              int      `json:"state_id"`
	CityID               int      `json:"city"`
	Sort                 string   `json:"sort"`
	Latitude             float32  `json:"lat"`
	Longitude            float32  `json:"lon"`
	Page                 int      `json:"page"`
	PerPage              int      `json:"per_page"`
	EquipmentCategories  []string `json:"equipment_categories"`
	SupplyCategories     []string `json:"supply_categories"`
	ContractPreferences  []int    `json:"contract_preferences"`
	RentTo               int64    `json:"rent_to"`
	RentFrom             int64    `json:"rent_from"`
	EquipmentSubcategory string   `json:"equipment_subcategory"`
}

func DecodeListQuotesRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	contractPreferences := make([]int, 0)
	filterContractPreferences, ok := r.URL.Query()["contract_preference[]"]
	if ok {
		for _, filterContractPreference := range filterContractPreferences {
			contractPreferences = append(contractPreferences, helpers.AlwaysIntFromString(filterContractPreference))
		}
	}
	return &ListQuotesRequest{
		Keywords:             helpers.SplitStringQuery(r.URL.Query().Get("keywords")),
		Zipcode:              r.URL.Query().Get("zipcode"),
		CityID:               helpers.AlwaysIntFromString(r.URL.Query().Get("city")),
		StateID:              helpers.AlwaysIntFromString(r.URL.Query().Get("state_id")),
		Sort:                 r.URL.Query().Get("sort"),
		Latitude:             helpers.AlwaysFloatFromString(r.URL.Query().Get("lat")),
		Longitude:            helpers.AlwaysFloatFromString(r.URL.Query().Get("lon")),
		Page:                 helpers.AlwaysIntFromString(r.URL.Query().Get("page")),
		PerPage:              helpers.AlwaysIntFromString(r.URL.Query().Get("per_page")),
		EquipmentCategories:  helpers.SplitStringQuery(r.URL.Query().Get("equipment_categories")),
		SupplyCategories:     helpers.SplitStringQuery(r.URL.Query().Get("supply_categories")),
		ContractPreferences:  contractPreferences,
		RentTo:               helpers.AlwaysInt64FromString(r.URL.Query().Get("rent_to")),
		RentFrom:             helpers.AlwaysInt64FromString(r.URL.Query().Get("rent_from")),
		EquipmentSubcategory: r.URL.Query().Get("equipment_subcategory"),
	}, nil
}

type GetByIDRequest struct {
	ID string `json:"id"`
}

func DecodeGetSingleQuoteRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	id, ok := mux.Vars(r)["id"]
	if !ok || id == "" {
		return nil, errors.New("please provide a valid quote id")
	}
	return &GetByIDRequest{ID: id}, nil
}
