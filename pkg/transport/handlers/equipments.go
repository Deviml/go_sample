package handlers

import (
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/equipments"
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	equipmentPath = "/equipments"
)

func MakeEquipmentHandler(r *mux.Router, ep endpoint.Endpoint) *mux.Router {
	r.
		PathPrefix(equipmentPath).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			equipments.DecodeGetEqRequest,
			http2.EncodeJSONResponse,
		))
	return r
}
