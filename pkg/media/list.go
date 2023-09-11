package media

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/go-kit/kit/log"
)

type ListAboutUsMediaRepository interface {
	ListAboutUs(ctx context.Context) ([]entities.Media, error)
}

type GetAboutUsService struct {
	logger     log.Logger
	r          ListAboutUsMediaRepository
	bucketName string
}

func NewGetAboutUsService(logger log.Logger, r ListAboutUsMediaRepository, bucketName string) *GetAboutUsService {
	return &GetAboutUsService{logger: logger, r: r, bucketName: bucketName}
}

func (g GetAboutUsService) GetAboutUsMedia(ctx context.Context) ([]entities.AboutUsMedia, error) {
	mediaList, err := g.r.ListAboutUs(ctx)
	if err != nil {
		return nil, err
	}
	aboutUsMediaList := []entities.AboutUsMedia{}
	for _, media := range mediaList {
		aboutUsMediaList = append(aboutUsMediaList, entities.AboutUsFromMedia(media, g.bucketName))
	}
	return aboutUsMediaList, nil
}
