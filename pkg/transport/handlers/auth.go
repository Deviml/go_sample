package handlers

import (
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/auth"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	authPrefix = "/auth"
	loginPath  = ""
)

func MakeAuthHandler(r *mux.Router, eps *auth.Endpoints) *mux.Router {
	addLoginRoute(r, eps.Login)
	return r
}

func addLoginRoute(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(authPrefix).
		Path(loginPath).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			requests.DecodeLoginRequest,
			http2.EncodeJSONResponse,
		))
}
