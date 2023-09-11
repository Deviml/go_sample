package requests

import (
	"context"
	"net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/helpers"
	"github.com/gorilla/mux"
)

type ListProposal struct {
	UserType string   `json:"type"`
	Keywords []string `json:"keywords"`
	Sort     string   `json:"sort"`
	Page     int      `json:"page"`
	PerPage  int      `json:"per_page"`
}

func DecodeListProposal(ctx context.Context, r *http.Request) (response interface{}, err error) {
	user_type, ok := mux.Vars(r)["type"]
	if !ok || user_type == "" {
		user_type = "buyer"
	}

	return &ListProposal{
		UserType: user_type,
		Keywords: helpers.SplitStringQuery(r.URL.Query().Get("keywords")),
		Sort:     r.URL.Query().Get("sort"),
		Page:     helpers.AlwaysIntFromString(r.URL.Query().Get("page")),
		PerPage:  helpers.AlwaysIntFromString(r.URL.Query().Get("per_page")),
	}, nil
}
