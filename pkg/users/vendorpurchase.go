package users

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/filters"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/sorts"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/sublists"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/vendorpurchases"
	"github.com/go-kit/kit/log"
	"time"
)

type VendorSublistsRepository interface {
	GetSublists(ctx context.Context, userID string, from string, to string) ([]entities.Sublist, error)
}

type GetSublistByIDRepository interface {
	GetSublistByID(ctx context.Context, sublistID string) (*models.Sublist, error)
}

type GetVendorSublists struct {
	vsr    sublists.ListSublistRepository
	logger log.Logger
	gsr    GetSublistByIDRepository
}

func NewGetVendorSublists(vsr sublists.ListSublistRepository, logger log.Logger, gsr GetSublistByIDRepository) *GetVendorSublists {
	return &GetVendorSublists{vsr: vsr, logger: logger, gsr: gsr}
}

func (g GetVendorSublists) GetByID(ctx context.Context, id string) (*entities.CompleteSublist, error) {
	sublist, err := g.gsr.GetSublistByID(ctx, id)
	if err != nil {
		return nil, err
	}
	sublistEntity := sublist.ToSublistEntity()
	return &sublistEntity, nil
}

func (g GetVendorSublists) GetSublists(ctx context.Context, request *vendorpurchases.GetVendorPurchasesRequest) ([]entities.Sublist, error) {
	from := ""
	if request.From != 0 {
		from = time.Unix(request.From, 0).Format(time.RFC3339)
	}
	to := ""
	if request.To != 0 {
		to = time.Unix(request.To, 0).Format(time.RFC3339)
	}
	s, err := g.vsr.FetchSublists(ctx, filters.ListSublists{
		Keywords: request.Keywords,
		From:     from,
		To:       to,
		UserID:   request.UserID,
	}, sorts.ListSublists{}, entities.PaginationQuery{})
	if err != nil {
		return nil, err
	}

	return makePublishUnix(s), nil
}

func makePublishUnix(sublists []entities.Sublist) []entities.Sublist {
	for idx, sublist := range sublists {
		sublist.PublishDateUnix = sublist.PurchasedAt.Unix()
		sublists[idx] = sublist
	}
	return sublists
}
