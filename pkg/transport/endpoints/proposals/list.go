package proposals

import (
	"context"
	"errors"
	"log"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/auth/jwt/claims"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

type ListProposalsService interface {
	List(ctx context.Context, userID string, paginationQuery entities.PaginationQuery, userType string) ([]entities.Proposal, error)
	GetPagination(ctx context.Context, userID string, paginationQuery entities.PaginationQuery, userType string) (*entities.Pagination, error)
}

func MakeListProposalsEndpoint(svc ListProposalsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		listRequest, ok := request.(*requests.ListProposal)
		var userID string
		token := ctx.Value(jwt.JWTClaimsContextKey)
		if token != nil {
			userToken := token.(*claims.UserClaims)
			userID = userToken.UserID
		}
		if !ok {
			return nil, errors.New("bad cast")
		}

		paginationQuery := entities.PaginationQuery{
			Page:    listRequest.Page,
			PerPage: listRequest.PerPage,
		}
		log.Printf(listRequest.UserType)
		proposals, err := svc.List(ctx, userID, paginationQuery, listRequest.UserType)
		if err != nil {
			return nil, err
		}

		pagination, err := svc.GetPagination(ctx, userID, paginationQuery, listRequest.UserType)
		if err != nil {
			return nil, err
		}

		return responses.CollectionResponse{
			Data:       proposals,
			Pagination: pagination,
		}, nil
	}
}
