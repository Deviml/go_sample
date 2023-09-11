package handlers

import (
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	equipmentCategoriesPrefix = "/equipment-categories"
	equipmentCategoriesList   = ""
)

func MakeEquipmentCategoriesHandler(r *mux.Router, ep endpoint.Endpoint) *mux.Router {
	addListEquipmentRoute(r, ep)
	return r
}

func addListEquipmentRoute(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(equipmentCategoriesPrefix).
		Path(equipmentCategoriesList).
		Methods(http.MethodGet).
		Handler(http2.NewServer(ep,
			http2.NopRequestDecoder,
			http2.EncodeJSONResponse,
		))
}
