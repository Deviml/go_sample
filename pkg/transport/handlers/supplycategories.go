package handlers

import (
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/supplycategories"
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	supplyCategoriesPrefix   = "/supply-categories"
	listSupplyCategoriesPath = ""
	supplies                 = "/supplies"
)

func MakeSupplyCategoriesHandler(r *mux.Router, ep endpoint.Endpoint) *mux.Router {
	addListSupplyCategoriesRoute(r, ep)
	return r
}

func MakeSuppliesHandler(r *mux.Router, ep endpoint.Endpoint) *mux.Router {
	r.
		Path(supplies).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			supplycategories.DecodeSupplies,
			http2.EncodeJSONResponse,
		))
	return r
}

func addListSupplyCategoriesRoute(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(supplyCategoriesPrefix).
		Path(listSupplyCategoriesPath).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			http2.NopRequestDecoder,
			http2.EncodeJSONResponse,
		))
}
