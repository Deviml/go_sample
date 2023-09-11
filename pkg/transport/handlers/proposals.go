package handlers

import (
	"net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/proposals"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const (
	proposalsPrefix   = "/proposals"
	proposalsTypePath = "/list/{type}"
	proposalsListPath = "/buyers"
	proposalsSellerListPath = "/sellers"
	createProposal    = "/store"
	updateProposal    = "/update/{id}"
	deleteProposal    = "/delete/{id}"
	showProposal      = "/show/{id}"
	acceptProposal    = "/accept/{id}"
	rejectProposal    = "/reject/{id}"
)

func MakeUserProposalHandler(r *mux.Router, eps proposals.Endpoints, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	addCreateProposalRoute(r, authMiddleware(eps.Create), op)
	addUpdateProposalRoute(r, authMiddleware(eps.Update), op)
	addShowProposalRoute(r, authMiddleware(eps.Show), op)
	addDeleteProposalRoute(r, authMiddleware(eps.Delete), op)

	return r
}

func MakeBuyerProposalHandler(r *mux.Router, eps proposals.Endpoints, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	addAcceptProposalRoute(r, authMiddleware(eps.Accept), op)
	addRejectProposalRoute(r, authMiddleware(eps.Reject), op)
	return r
}

func MakeBuyProposalListHandler(r *mux.Router, ep endpoint.Endpoint, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	r.
		PathPrefix(proposalsPrefix).
		PathPrefix(proposalsListPath).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			authMiddleware(ep),
			proposals.DecodeGetBuyerProposalRequest,
			http2.EncodeJSONResponse,
			op...,
		))
	return r
}

func MakeSellerProposalListHandler(r *mux.Router, ep endpoint.Endpoint, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	r.
		PathPrefix(proposalsPrefix).
		PathPrefix(proposalsSellerListPath).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			authMiddleware(ep),
			proposals.DecodeGetSellerProposalRequest,
			http2.EncodeJSONResponse,
			op...,
		))
	return r
}
func addCreateProposalRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(proposalsPrefix).
		PathPrefix(createProposal).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			proposals.DecodeCreateRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addUpdateProposalRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(proposalsPrefix).
		PathPrefix(updateProposal).
		Methods(http.MethodPatch).
		Handler(http2.NewServer(
			ep,
			proposals.DecodeUpdateRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addShowProposalRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(proposalsPrefix).
		PathPrefix(showProposal).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			proposals.DecodeShowProposalRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addDeleteProposalRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(proposalsPrefix).
		PathPrefix(deleteProposal).
		Methods(http.MethodDelete).
		Handler(http2.NewServer(
			ep,
			proposals.DecodeDeleteRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addAcceptProposalRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(proposalsPrefix).
		PathPrefix(acceptProposal).
		Methods(http.MethodPatch).
		Handler(http2.NewServer(
			ep,
			proposals.DecodeDeleteRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addRejectProposalRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(proposalsPrefix).
		PathPrefix(rejectProposal).
		Methods(http.MethodPatch).
		Handler(http2.NewServer(
			ep,
			proposals.DecodeDeleteRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func MakeProposalListHandler(r *mux.Router, ep endpoint.Endpoint, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	addListProposalRoute(r, ep, authMiddleware, op)
	return r
}

func addListProposalRoute(r *mux.Router, ep endpoint.Endpoint, authMiddleware endpoint.Middleware, op []http2.ServerOption) {
	r.PathPrefix(proposalsPrefix).
		Path(proposalsTypePath).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			authMiddleware(ep),
			requests.DecodeListProposal,
			http2.EncodeJSONResponse,
			op...,
		))
}