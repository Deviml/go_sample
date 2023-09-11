package handlers

import (
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/companies"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	companyPrefix     = "/companies"
	listCompaniesPath = ""
)

func MakeCompaniesHandler(r *mux.Router, eps *companies.Endpoints) *mux.Router {
	addListCompaniesRoute(r, eps.List)
	return r
}

func addListCompaniesRoute(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(companyPrefix).
		Path(listCompaniesPath).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			requests.DecodeListCompaniesRequest,
			http2.EncodeJSONResponse,
		))
}
