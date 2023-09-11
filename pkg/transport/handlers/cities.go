package handlers

import (
	"net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/cities"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const (
	citiesPrefix = "/cities"
	getCities    = ""
	counties     = "/counties"
	getCounties  = ""
	getStates    = "/states"
)

func MakeCitiesHandler(r *mux.Router, eps *cities.Endpoints) *mux.Router {
	r.PathPrefix(citiesPrefix).
		Path(getCities).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			eps.GetCities,
			cities.DecodeGetCitiesRequest,
			http2.EncodeJSONResponse,
		))

	r.PathPrefix(counties).
		Path(getCounties).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			eps.GetCounties,
			cities.DecodeGetCitiesRequest,
			http2.EncodeJSONResponse,
		))
	r.Path(getStates).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			eps.GetState,
			cities.DecodeGetCitiesRequest,
			http2.EncodeJSONResponse,
		))
	return r
}
