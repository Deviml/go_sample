package requests

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/helpers"
	"net/http"
)

type GetSublistsRequest struct {
	Keywords  []string `json:"keywords"`
	Zipcode   string   `json:"zipcode"`
	StateID   int      `json:"state"`
	CityID    int      `json:"city"`
	Sort      string   `json:"sort"`
	Latitude  float32  `json:"lat"`
	Longitude float32  `json:"lon"`
	Page      int      `json:"page"`
	PerPage   int      `json:"per_page"`
}

func DecodeGetSublistRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	return &GetSublistsRequest{
		Keywords:  helpers.SplitStringQuery(r.URL.Query().Get("keywords")),
		Zipcode:   r.URL.Query().Get("zipcode"),
		CityID:    helpers.AlwaysIntFromString(r.URL.Query().Get("city")),
		StateID:   helpers.AlwaysIntFromString(r.URL.Query().Get("city")),
		Sort:      r.URL.Query().Get("sort"),
		Latitude:  helpers.AlwaysFloatFromString(r.URL.Query().Get("lat")),
		Longitude: helpers.AlwaysFloatFromString(r.URL.Query().Get("lon")),
		Page:      helpers.AlwaysIntFromString(r.URL.Query().Get("page")),
		PerPage:   helpers.AlwaysIntFromString(r.URL.Query().Get("per_page")),
	}, nil
}
