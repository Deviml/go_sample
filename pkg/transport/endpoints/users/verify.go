package users

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/auth/jwt/claims"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

type VerifyService interface {
	Verify(ctx context.Context, userID string, code string, codeType int) error
}

func MakeVerifySMSEndpoint(svc VerifyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		verifyRequest := request.(*requests.VerifyRequest)
		err = svc.Verify(ctx, userID, verifyRequest.Code, models.SMSVerificationType)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}

func MakeVerifyEmailEndpoint(svc VerifyService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		verifyRequest := request.(*requests.VerifyRequest)
		err = svc.Verify(ctx, userID, verifyRequest.Code, models.EmailVerificationType)
		if err != nil {
			return nil, err
		}
		return responses.NoContent{}, nil
	}
}
