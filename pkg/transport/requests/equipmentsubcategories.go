package requests

import (
	"context"
	"net/http"
)

type GetEquipmentSubcategories struct {
	CategoryID string
}

func DecodeGetEquipmentSubcategories(ctx context.Context, r *http.Request) (response interface{}, err error) {
	category := r.URL.Query().Get("category")
	if category == "" {
		category = r.URL.Query().Get("equipment_category")
	}
	return &GetEquipmentSubcategories{
		CategoryID: category,
	}, nil
}
