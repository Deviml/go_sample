package handlers

import (
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	mediaPrefix = "/media"
	aboutUsPath = "/about-us"
)

func MakeMediaHandler(r *mux.Router, ep endpoint.Endpoint) *mux.Router {
	addListAboutUsRoute(r, ep)
	return r
}

func addListAboutUsRoute(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(mediaPrefix).
		Path(aboutUsPath).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			http2.NopRequestDecoder,
			http2.EncodeJSONResponse,
		))
}
