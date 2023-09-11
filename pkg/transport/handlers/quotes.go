package handlers

import (
	"net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/quotes"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const (
	quotesPrefix    = "/quotes"
	quotesListPath  = ""
	createEquipment = "/equipment-request"
	updateEquipment = "/equipment-requests/{id}"
	deleteEquipment = "/equipment-requests/{id}"
	showEquipment   = "/equipment-requests/{id}"
	createSupply    = "/supply-requests"
	updateSupply    = "/supply-requests/{id}"
	deleteSupply    = "/supply-requests/{id}"
	showSupply      = "/supply-requests/{id}"
	accountsQuote   = "/accounts/subcontractor/quotes"
)

func MakeQuotesHandler(r *mux.Router, ep endpoint.Endpoint) *mux.Router {
	addListQuotesRoute(r, ep)
	return r
}

// Reference from this TODO
func addListQuotesRoute(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(quotesPrefix).
		Path(quotesListPath).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			requests.DecodeListQuotesRequest,
			http2.EncodeJSONResponse,
		))
}

// Working on this TODO
func MakeUserQuotesHandler(r *mux.Router, ep endpoint.Endpoint, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	r.
		PathPrefix(accountsQuote).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			authMiddleware(ep),
			quotes.DecodeGetUserQuotesRequest,
			http2.EncodeJSONResponse,
			op...,
		))
	return r
}

func MakeUserSupplyHandler(r *mux.Router, eps quotes.Endpoints, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	addCreateSupplyRoute(r, authMiddleware(eps.CreateSupply), op)
	addUpdateSupplyRoute(r, authMiddleware(eps.UpdateSupply), op)
	addShowSupplyRoute(r, authMiddleware(eps.ShowSupply), op)
	addDeleteSupplyRoute(r, authMiddleware(eps.DeleteSupply), op)
	return r
}

func addDeleteSupplyRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(quotesPrefix).
		PathPrefix(deleteSupply).
		Methods(http.MethodDelete).
		Handler(http2.NewServer(
			ep,
			quotes.DecodeDeleteSupplyRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addShowSupplyRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(quotesPrefix).
		PathPrefix(showSupply).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			quotes.DecodeShowSupplyRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addUpdateSupplyRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(quotesPrefix).
		PathPrefix(updateSupply).
		Methods(http.MethodPatch).
		Handler(http2.NewServer(
			ep,
			quotes.DecodeUpdateSupplyRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addCreateSupplyRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(quotesPrefix).
		PathPrefix(createSupply).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			quotes.DecodeCreateSupplyRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func MakeUserEquipmentHandler(r *mux.Router, eps quotes.Endpoints, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	addQuoteCreateEquipmentRoute(r, authMiddleware(eps.CreateEquipment), op)
	addUpdateEquipmentRoute(r, authMiddleware(eps.UpdateEquipment), op)
	addShowEquipmentRoute(r, authMiddleware(eps.ShowEquipment), op)
	addDeleteEquipmentRoute(r, authMiddleware(eps.DeleteEquipment), op)
	return r
}

func addQuoteCreateEquipmentRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(quotesPrefix).
		PathPrefix(createEquipment).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			quotes.DecodeCreateEquipmentRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addUpdateEquipmentRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(quotesPrefix).
		PathPrefix(updateEquipment).
		Methods(http.MethodPatch).
		Handler(http2.NewServer(
			ep,
			quotes.DecodeUpdateEquipmentRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addShowEquipmentRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(quotesPrefix).
		PathPrefix(showEquipment).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			quotes.DecodeShowEquipmentRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addDeleteEquipmentRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(quotesPrefix).
		PathPrefix(deleteEquipment).
		Methods(http.MethodDelete).
		Handler(http2.NewServer(
			ep,
			quotes.DecodeDeleteEquipmentRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}
