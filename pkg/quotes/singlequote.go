package quotes

import (
	"context"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/go-kit/kit/log"
)

type SingleQuoteRepository interface {
	GetByID(ctx context.Context, id string) (entities.CompleteSingleQuote, error)
}

type SingleQuoteService struct {
	logger log.Logger
	qr     SingleQuoteRepository
}

func NewSingleQuoteService(logger log.Logger, qr SingleQuoteRepository) *SingleQuoteService {
	return &SingleQuoteService{logger: logger, qr: qr}
}

func (s SingleQuoteService) GetByID(ctx context.Context, id string) (entities.CompleteSingleQuote, error) {
	return s.qr.GetByID(ctx, id)
}
