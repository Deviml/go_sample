package media

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/endpoint"
)

type GetAboutUsMediaService interface {
	GetAboutUsMedia(ctx context.Context) ([]entities.AboutUsMedia, error)
}

func MakeAboutUsEndpoint(svc GetAboutUsMediaService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		aboutUsMedia, err := svc.GetAboutUsMedia(ctx)
		return responses.CollectionResponse{
			Data:       aboutUsMedia,
			Pagination: nil,
		}, nil
	}
}
