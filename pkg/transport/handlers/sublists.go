package handlers

import (
	sublists2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/sublists"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	sublistsPrefix    = "/sublists"
	getSublistsPath   = ""
	createSublist     = ""
	showSublist       = "/{id}"
	pathSublist       = "/{id}"
	deleteSublist     = "/{id}"
	userSublist       = "/accounts/general-contractor/sublists"
	userSublistDetail = "/accounts/general-contractor/sublists/{id}"
)

func MakeUserSubHandler(r *mux.Router, ep endpoint.Endpoint, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	r.
		Path(userSublist).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			authMiddleware(ep),
			http2.NopRequestDecoder,
			http2.EncodeJSONResponse,
			op...,
		))
	return r
}

func MakeGetSublistHandler(r *mux.Router, ep endpoint.Endpoint) *mux.Router {
	r.PathPrefix(sublistsPrefix).
		Path(getSublistsPath).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			requests.DecodeGetSublistRequest,
			http2.EncodeJSONResponse,
		))
	return r
}

func MakeUserSublistHandler(r *mux.Router, eps *sublists2.Endpoints, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	addCreateSublistRoute(r, authMiddleware(eps.Create), op)
	addShowSublistroute(r, authMiddleware(eps.Show), op)
	addUpdateSublistRoute(r, authMiddleware(eps.Update), op)
	addDeleteSublistRoute(r, authMiddleware(eps.Delete), op)
	return r
}

func addCreateSublistRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(sublistsPrefix).
		Path(createSublist).
		Methods(http.MethodPost).
		Handler(http2.NewServer(ep, sublists2.DecodeCreateRequest, http2.EncodeJSONResponse, op...))
}

func addShowSublistroute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(sublistsPrefix).
		Path(showSublist).
		Methods(http.MethodGet).
		Handler(http2.NewServer(ep, sublists2.DecodeShowRequest, http2.EncodeJSONResponse, op...))
	r.Path(userSublistDetail).Methods(http.MethodGet).Handler(http2.NewServer(ep, sublists2.DecodeShowRequest, http2.EncodeJSONResponse, op...))
}

func addUpdateSublistRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(sublistsPrefix).
		Path(pathSublist).
		Methods(http.MethodPatch).
		Handler(http2.NewServer(ep, sublists2.DecodePatchRequest, http2.EncodeJSONResponse, op...))
}

func addDeleteSublistRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(sublistsPrefix).
		Path(deleteSublist).
		Methods(http.MethodDelete).
		Handler(http2.NewServer(ep, sublists2.DecodeDeleteRequest, http2.EncodeJSONResponse, op...))
}
