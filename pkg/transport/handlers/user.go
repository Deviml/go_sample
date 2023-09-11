package handlers

import (
	"net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/users"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const (
	accountsPrefix    = "/accounts"
	vendorRentalPath  = "/vendor-rental"
	subcontractorPath = "/subcontractor"
	generalContractor = "/general-contractor"
	verifySMS         = "/verify-sms"
	verifyEmail       = "/verify-email"
	resendCode        = "/resend-code"
	closeAccount      = "/close-account"
	updateAccount     = "/update-account"
	profile           = "/me"
	requestInvitation = "/general-contractor/request-invitation"
	forgotPassword    = "/forgot-password"
)

func MakeUserHandler(r *mux.Router, eps *users.Endpoints, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	addCreateVendorRentalRoute(r, eps.CreateVendorRental)
	addCreateSubcontractorRoute(r, eps.CreateSubcontractor)
	addCreateGeneralContractorRoute(r, eps.CreateGeneralContractor)
	addForgotPassword(r, eps.ForgotPassword)
	addRequestInvitationCode(r, eps.RequestInvitationCode)
	addVerifySMSRoute(r, authMiddleware(eps.VerifySMS), op)
	addVerifyEmailRoute(r, authMiddleware(eps.VerifyEmail), op)
	addResendCodeRoute(r, authMiddleware(eps.ResendCode), op)
	addCloseAccountRoute(r, authMiddleware(eps.CloseAccount), op)
	addUpdateAccountRoute(r, authMiddleware(eps.UpdateAccount), op)
	addGetUserRoute(r, authMiddleware(eps.GetUser), op)
	addPatchUserRoute(r, authMiddleware(eps.PatchUser), op)
	return r
}

func addForgotPassword(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(accountsPrefix).
		Path(forgotPassword).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			requests.DecodeForgotPassword,
			http2.EncodeJSONResponse,
		))
}

func addCloseAccountRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(closeAccount).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			http2.NopRequestDecoder,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addUpdateAccountRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(updateAccount).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			requests.DecodeUpdateAccountRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addRequestInvitationCode(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(accountsPrefix).
		Path(requestInvitation).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			requests.DecodeRequestInvitationRequest,
			http2.EncodeJSONResponse,
		))
}

func addPatchUserRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(profile).
		Methods(http.MethodPatch).
		Handler(http2.NewServer(
			ep,
			requests.DecodeUpdateUserRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addGetUserRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(profile).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			http2.NopRequestDecoder,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addVerifyEmailRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(verifyEmail).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			requests.DecodeVerifyRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addVerifySMSRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(verifySMS).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			requests.DecodeVerifyRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addResendCodeRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(resendCode).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			requests.DecodeCreateResendCodeRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addCreateVendorRentalRoute(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(accountsPrefix).
		Path(vendorRentalPath).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			requests.DecodeCreateVendorRentalRequest,
			http2.EncodeJSONResponse,
		))
}

func addCreateSubcontractorRoute(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(accountsPrefix).
		Path(subcontractorPath).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			requests.DecodeCreateSubcontractorRequest,
			http2.EncodeJSONResponse,
		))
}

func addCreateGeneralContractorRoute(r *mux.Router, ep endpoint.Endpoint) {
	r.PathPrefix(accountsPrefix).
		Path(generalContractor).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			requests.DecodeCreateGeneralContractorRequest,
			http2.EncodeJSONResponse,
		))
}
