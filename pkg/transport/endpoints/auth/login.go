package auth

import (
	"context"
	"strconv"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/errors/http"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/endpoint"
)

type LoginService interface {
	Login(ctx context.Context, username string, password string) (*entities.LoginInformation, error)
}

func MakeLoginEndpoint(svc LoginService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		loginRequest, _ := request.(*requests.LoginRequest)

		claims, err := svc.Login(ctx, loginRequest.Username, loginRequest.Password)
		if err != nil {
			return nil, http.InvalidLoginError{}
		}
		response, err = LoginInformationToLoginResponse(claims)

		if err != nil {
			return nil, http.BadRequest{}
		}
		return response, nil
	}
}

func LoginInformationToLoginResponse(loginInformation *entities.LoginInformation) (*responses.LoginResponse, error) {
	return &responses.LoginResponse{
		Token:        loginInformation.Token,
		RefreshToken: loginInformation.RefreshToken,
		Exp:          strconv.Itoa(int(loginInformation.Claims.ExpiresAt)),
		UserID: 	 loginInformation.UserID,
	}, nil
}
