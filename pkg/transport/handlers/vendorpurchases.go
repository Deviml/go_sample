package handlers

import (
	"net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/vendorpurchases"
	"github.com/go-kit/kit/endpoint"
	http2 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const (
	supplyRequests                        = "/supply-requests"
	supplyDetail                          = "/supply-requests/{id}"
	equipmentRequestsRD                   = "/equipment-requests"
	equipmentRequestsRA                   = "/equipment-requests/sorted"
	equipmentRequestsPD                   = "/equipment-requests/sorted-pd"
	equipmentRequestsPA                   = "/equipment-requests/sorted-pa"
	equipmentDetail                       = "/equipment-requests/{id}"
	sublists                              = "/sublists"
	sublistsDetail                        = "/sublists/{id}"
	checkoutPrefix                        = "/checkout"
	checkout                              = ""
	checkoutwithoutpaymentPrefix          = "/checkout-discount"
	checkoutwithoutpayment                = ""
	paymentIntentPrefix                   = "/paymentIntent"
	paymentIntent                         = ""
	ConfirmPaymentIntentAndCheckoutPrefix = "/confirmPaymentIntentAndCheckout"
	ConfirmPaymentIntentAndCheckout       = ""
)

func MakeVendorPurchasesHandler(r *mux.Router, eps *vendorpurchases.Endpoints, authMiddleware endpoint.Middleware, op []http2.ServerOption) *mux.Router {
	addAccountSupplyRequestRoute(r, authMiddleware(eps.GetSupplyRequests), op)
	addAccountEquipmentRequestSortRDRoute(r, authMiddleware(eps.GetEquipmentRequestsSortedRD), op)
	addAccountEquipmentRequestSortRARoute(r, authMiddleware(eps.GetEquipmentRequestsSortedRA), op)
	addAccountEquipmentRequestSortPARoute(r, authMiddleware(eps.GetEquipmentRequestsSortedPA), op)
	addAccountEquipmentRequestSortPDRoute(r, authMiddleware(eps.GetEquipmentRequestsSortedPD), op)
	addAccountSublists(r, authMiddleware(eps.GetSublists), op)
	addAccountSublistsDetail(r, authMiddleware(eps.GetSublistsDetail), op)
	addAccountSupplyRequestDetailroute(r, authMiddleware(eps.GetSupplyDetail), op)
	addAccountEquipmentRequestDetailRoute(r, authMiddleware(eps.GetEquipmentDetail), op)
	addCheckoutEndpoint(r, authMiddleware(eps.Checkout), op)
	addPaymentIntentEndpoint(r, authMiddleware(eps.PaymentIntent), op)
	addConfirmPaymentIntentAndCheckoutEndpoint(r, authMiddleware(eps.ConfirmPaymentIntentAndCheckout), op)
	return r
}

func addCheckoutEndpoint(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(checkoutPrefix).
		Path(checkout).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodeCheckoutRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addPaymentIntentEndpoint(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(paymentIntentPrefix).
		Path(paymentIntent).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodePaymentIntentRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addConfirmPaymentIntentAndCheckoutEndpoint(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(ConfirmPaymentIntentAndCheckoutPrefix).
		Path(ConfirmPaymentIntentAndCheckout).
		Methods(http.MethodPost).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodeCheckoutRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addAccountEquipmentRequestDetailRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(equipmentDetail).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodePurchaseDetailRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addAccountSupplyRequestDetailroute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(supplyDetail).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodePurchaseDetailRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addAccountSublistsDetail(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(sublistsDetail).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodePurchaseDetailRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addAccountEquipmentRequestSortRARoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(equipmentRequestsRA).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodeGetVendorSublistsRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addAccountEquipmentRequestSortRDRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(equipmentRequestsRD).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodeGetVendorSublistsRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addAccountEquipmentRequestSortPARoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(equipmentRequestsPA).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodeGetVendorSublistsRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addAccountEquipmentRequestSortPDRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(equipmentRequestsPD).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodeGetVendorSublistsRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addAccountSublists(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(sublists).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodeGetVendorSublistsRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}

func addAccountSupplyRequestRoute(r *mux.Router, ep endpoint.Endpoint, op []http2.ServerOption) {
	r.PathPrefix(accountsPrefix).
		Path(supplyRequests).
		Methods(http.MethodGet).
		Handler(http2.NewServer(
			ep,
			vendorpurchases.DecodeGetVendorSublistsRequest,
			http2.EncodeJSONResponse,
			op...,
		))
}
