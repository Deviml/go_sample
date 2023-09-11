package handlers

import (
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	equipmentSubcategoriesPrefix = "/equipment-subcategories"
	listEquipmentSubcategories   = ""
)

func MakeEquipmentSubcategoriesHandler(r *mux.Router, ep endpoint.Endpoint) *mux.Router {
	addListEquipmentSubcategoriesRoute(r, ep)
	return r
}

func addListEquipmentSubcategoriesRoute(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(equipmentSubcategoriesPrefix).
		Path(listEquipmentSubcategories).
		Methods(http.MethodGet).
		Handler(http2.NewServer(ep,
			requests.DecodeGetEquipmentSubcategories,
			http2.EncodeJSONResponse,
		))
}
