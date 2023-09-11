package users

import (
	"context"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/auth/jwt/claims"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/errors/http"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/auth"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

type Services struct {
	CreateUserService
	VerifyService
	GetUserService
	PatchUserService
}

func NewServices(createUserService CreateUserService, verifyService VerifyService, getUserService GetUserService, patchUserService PatchUserService) *Services {
	return &Services{CreateUserService: createUserService, VerifyService: verifyService, GetUserService: getUserService, PatchUserService: patchUserService}
}

type Endpoints struct {
	CreateVendorRental      endpoint.Endpoint
	CreateSubcontractor     endpoint.Endpoint
	CreateGeneralContractor endpoint.Endpoint
	VerifySMS               endpoint.Endpoint
	VerifyEmail             endpoint.Endpoint
	ResendCode              endpoint.Endpoint
	GetUser                 endpoint.Endpoint
	PatchUser               endpoint.Endpoint
	RequestInvitationCode   endpoint.Endpoint
	ForgotPassword          endpoint.Endpoint
	CloseAccount            endpoint.Endpoint
	UpdateAccount           endpoint.Endpoint
}

func NewEndpoints(svc *Services) *Endpoints {
	return &Endpoints{
		CreateVendorRental:      makeCreateVendorRentalEndpoint(svc),
		CreateSubcontractor:     makeCreateSubcontractorEndpoint(svc),
		CreateGeneralContractor: makeCreateGeneralContractor(svc),
		VerifySMS:               MakeVerifySMSEndpoint(svc),
		VerifyEmail:             MakeVerifyEmailEndpoint(svc),
		ResendCode:              MakeResendCodeEndpoint(svc),
		GetUser:                 makeGetUserEndpoint(svc),
		PatchUser:               makePatchUserEndpoint(svc),
		RequestInvitationCode:   makeRequestInvitationCode(svc),
		ForgotPassword:          makeForgotPassword(svc),
		CloseAccount:            makeCloseAccountEndpoint(svc),
		UpdateAccount:           makeUpdateAccountEndpoint(svc),
	}
}

type PatchUserService interface {
	PatchUser(ctx context.Context, userID string, request *requests.UpdateUserRequest) error
}

func makePatchUserEndpoint(svc PatchUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		updateRequest := request.(*requests.UpdateUserRequest)
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token == nil {
			return nil, http.BadRequest{}
		}

		userToken, ok := token.(*claims.UserClaims)
		if !ok {
			return nil, http.BadRequest{}
		}
		err = svc.PatchUser(ctx, userToken.UserID, updateRequest)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func makeRequestInvitationCode(svc CreateUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ri := request.(*requests.RequestInvitation)

		err = svc.RequestInvitation(ctx, ri)

		if err != nil {
			return nil, http.BadRequest{}
		}
		return responses.NoContent{}, nil
	}
}

type GetUserService interface {
	GetUser(ctx context.Context, userID string) (*entities.UserInfo, error)
	CloseAccount(ctx context.Context, userID string) error
	UpdateAccount(ctx context.Context, userID string, req *requests.UpdateAccount) error
}

func makeCloseAccountEndpoint(svc GetUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token == nil {
			return nil, http.BadRequest{}
		}

		userToken, ok := token.(*claims.UserClaims)
		if !ok {
			return nil, http.BadRequest{}
		}
		err = svc.CloseAccount(ctx, userToken.UserID)
		if err != nil {
			return nil, err
		}
		return responses.SuccessMessage{Message: "Account closed successfully!"}, nil
	}
}

func makeUpdateAccountEndpoint(svc GetUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*requests.UpdateAccount)
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token == nil {
			return nil, http.BadRequest{}
		}
		userToken, ok := token.(*claims.UserClaims)
		if !ok {
			return nil, http.BadRequest{}
		}
		err = svc.UpdateAccount(ctx, userToken.UserID, req)
		if err != nil {
			return nil, err
		}
		return responses.SuccessMessage{Message: "Account updated successfully!"}, nil
	}
}

func makeGetUserEndpoint(svc GetUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token == nil {
			return nil, http.BadRequest{}
		}

		userToken, ok := token.(*claims.UserClaims)
		if !ok {
			return nil, http.BadRequest{}
		}
		user, err := svc.GetUser(ctx, userToken.UserID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{
			Data: user,
		}, nil
	}
}

type CreateUserService interface {
	CreateVendorRental(ctx context.Context, user entities.User, profile entities.VendorRental, confirmationType string) (*entities.LoginInformation, error)
	CreateSubcontractor(ctx context.Context, user entities.User, profile entities.Subcontractor, confirmationType string) (*entities.LoginInformation, error)
	CreateGeneralContractor(ctx context.Context, user entities.User, profile entities.GeneralContractor, confirmationType string, verificationCode string) (*entities.LoginInformation, error)
	RequestInvitation(ctx context.Context, ri *requests.RequestInvitation) error
	ForgotPassword(ctx context.Context, fp *requests.ForgotPassword) error
	ResendCode(ctx context.Context, Email string, confirmationType string) error
}

func MakeResendCodeEndpoint(svc CreateUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*requests.CreateResendCodeRequest)
		err = svc.ResendCode(ctx, req.UserID, req.ConfirmationMethod)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func makeCreateVendorRentalEndpoint(svc CreateUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		createRequest := request.(*requests.CreateVendorRentalRequest)
		claims, err := svc.CreateVendorRental(ctx, entities.User{
			Password:    createRequest.Password,
			Email:       createRequest.Email,
			Phone:       createRequest.Phone,
			FullName:    createRequest.FullName,
			CityName:    createRequest.CityName,
			StateName:   createRequest.StateName,
			CountryName: createRequest.CountryName,
			Address:     createRequest.Address,
			Zipcode:     createRequest.Zipcode,
		}, entities.VendorRental{
			Profile: entities.Profile{
				Type: entities.VendorRentalProfileType,
			},
			StateID: createRequest.StateID,
			CityID:  createRequest.CityID,
			Company: entities.Company{
				Name: createRequest.CompanyName,
			},
		}, createRequest.ConfirmationMethod)
		if err != nil {
			return nil, http.ExistingCredentials{}
		}
		response, err = auth.LoginInformationToLoginResponse(claims)

		if err != nil {
			return nil, http.BadRequest{}
		}
		return response, nil
	}
}

func makeCreateSubcontractorEndpoint(svc CreateUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		createRequest := request.(*requests.CreateSubcontractorRequest)
		claims, err := svc.CreateSubcontractor(ctx, entities.User{
			Password:    createRequest.Password,
			Email:       createRequest.Email,
			Phone:       createRequest.Phone,
			FullName:    createRequest.FullName,
			CityName:    createRequest.CityName,
			StateName:   createRequest.StateName,
			CountryName: createRequest.CountryName,
			Address:     createRequest.Address,
			Zipcode:     createRequest.Zipcode,
		}, entities.Subcontractor{
			Profile: entities.Profile{
				Type: entities.SubcontractorProfileType,
			},
			StateID: createRequest.StateID,
			CityID:  createRequest.CityID,
			Company: entities.Company{
				Name: createRequest.CompanyName,
			},
		}, createRequest.ConfirmationMethod)
		if err != nil {
			return nil, http.InvalidLoginError{}
		}
		response, err = auth.LoginInformationToLoginResponse(claims)

		if err != nil {
			return nil, http.BadRequest{}
		}
		return response, nil
	}
}

func makeCreateGeneralContractor(svc CreateUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		createRequest := request.(*requests.CreateGeneralContractorRequest)
		claims, err := svc.CreateGeneralContractor(ctx, entities.User{
			Password:    createRequest.Password,
			Email:       createRequest.Email,
			Phone:       createRequest.Phone,
			FullName:    createRequest.FullName,
			CityName:    createRequest.CityName,
			StateName:   createRequest.StateName,
			CountryName: createRequest.CountryName,
			Address:     createRequest.Address,
			Zipcode:     createRequest.Zipcode,
		}, entities.GeneralContractor{
			Profile: entities.Profile{
				Type: entities.GeneralContractorProfileType,
			},
			StateID: createRequest.StateID,
			CityID:  createRequest.CityID,
			Company: entities.Company{
				Name: createRequest.CompanyName,
			},
			VerificationCode: createRequest.VerificationCode,
		}, createRequest.ConfirmationMethod, createRequest.VerificationCode)
		if err != nil {
			return nil, http.InvalidLoginError{}
		}
		response, err = auth.LoginInformationToLoginResponse(claims)

		if err != nil {
			return nil, http.BadRequest{}
		}
		return response, nil
	}
}

func makeForgotPassword(svc CreateUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ri := request.(*requests.ForgotPassword)

		err = svc.ForgotPassword(ctx, ri)

		if err != nil {
			return nil, http.BadRequest{}
		}
		return responses.NoContent{}, nil
	}
}
