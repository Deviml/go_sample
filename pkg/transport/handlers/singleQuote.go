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
	SingleQuotePrefix = "/quote"
	GetSingleQuote    = "/{id}"
)

func MakeSingleQuoteHandler(r *mux.Router, ep *quotes.Endpoint) *mux.Router {
	addGetSingleQuoteRoute(r, ep.GetByID)
	return r
}

func addGetSingleQuoteRoute(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(SingleQuotePrefix).
		Path(GetSingleQuote).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			requests.DecodeGetSingleQuoteRequest,
			http2.EncodeJSONResponse,
		))
}