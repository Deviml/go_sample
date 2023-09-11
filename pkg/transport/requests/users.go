package requests

import (
	"context"
	"encoding/json"
	"net/http"

	http2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/errors/http"
)

type CreateVendorRentalRequest struct {
	StateID            string `json:"state_id"`
	CityID             string `json:"city_id"`
	CityName           string `json:"city_Name"`
	StateName          string `json:"state_Name"`
	CountryName        string `json:"countryName"`
	Address            string `json:"address"`
	Zipcode            string `json:"zipcode"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	Phone              string `json:"phone"`
	FullName           string `json:"full_name"`
	CompanyName        string `json:"company_name"`
	ConfirmationMethod string `json:"confirmation_method"`
}

func DecodeCreateVendorRentalRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var create CreateVendorRentalRequest
	err = json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &create, nil
}

type CreateResendCodeRequest struct {
	UserID             string `json:"id"`
	ConfirmationMethod string `json:"confirmation_method"`
}

func DecodeCreateResendCodeRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var create CreateResendCodeRequest
	err = json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &create, nil
}

type CreateSubcontractorRequest struct {
	StateID            string `json:"state_id"`
	CityID             string `json:"city_id"`
	CityName           string `json:"city_Name"`
	StateName          string `json:"state_Name"`
	CountryName        string `json:"countryName"`
	Address            string `json:"address"`
	Zipcode            string `json:"zipcode"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	Phone              string `json:"phone"`
	FullName           string `json:"full_name"`
	CompanyName        string `json:"company_name"`
	ConfirmationMethod string `json:"confirmation_method"`
}

func DecodeCreateSubcontractorRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var create CreateSubcontractorRequest
	err = json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &create, nil
}

type CreateGeneralContractorRequest struct {
	StateID            string `json:"state_id"`
	CityID             string `json:"city_id"`
	CityName           string `json:"city"`
	StateName          string `json:"state"`
	CountryName        string `json:"country"`
	Address            string `json:"address"`
	Zipcode            string `json:"zipcode"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	Phone              string `json:"phone"`
	FullName           string `json:"full_name"`
	CompanyName        string `json:"company_name"`
	ConfirmationMethod string `json:"confirmation_method"`
	VerificationCode   string `json:"invitationCode"`
}

func DecodeCreateGeneralContractorRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var create CreateGeneralContractorRequest
	err = json.NewDecoder(r.Body).Decode(&create)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &create, nil
}

type UpdateUserRequest struct {
	SublistNotification          bool   `json:"sublist_notifications"`
	EquipmentRequestNotification bool   `json:"equipment_notifications"`
	SupplyRequestNotification    bool   `json:"supplies_notifications"`
	EveryStateSelection          bool   `json:"every_state_selection"`
	NewUserPopUp                 bool   `json:"new_user_pop_up"`
	SupplyCategories             []uint `json:"supply_categories"`
	EquipmentCategories          []uint `json:"equipment_categories"`
	Cities                       []uint `json:"cities"`
	States                       []uint `json:"states"`
	Counties                     []uint `json:"counties"`
	Password                     string `json:"password"`
}

func DecodeUpdateUserRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var udpate UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&udpate)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &udpate, nil
}

type RequestInvitation struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	Name  string `json:"full_name"`
}

func DecodeRequestInvitationRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var ri RequestInvitation
	err = json.NewDecoder(r.Body).Decode(&ri)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &ri, nil
}

type ForgotPassword struct {
	Email string `json:"email"`
}

func DecodeForgotPassword(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var fp ForgotPassword
	err = json.NewDecoder(r.Body).Decode(&fp)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &fp, nil
}

type UpdateAccount struct {
	FullName       string `json:"full_name"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profile_picture"`
}

func DecodeUpdateAccountRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var ua UpdateAccount
	err = json.NewDecoder(r.Body).Decode(&ua)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &ua, nil
}
