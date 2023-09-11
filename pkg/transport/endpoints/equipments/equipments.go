package equipments

import (
	"context"
	"net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/responses"
	"github.com/go-kit/kit/endpoint"
)

type GetEquipmentService interface {
	Get(ctx context.Context, subID string) ([]entities.Equipment, error)
}

type GetEqRequest struct {
	SubID string
}

func DecodeGetEqRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return &GetEqRequest{SubID: r.URL.Query().Get("subcategory")}, nil
}

func MakeEquipmentEndpoint(svc GetEquipmentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		gr := request.(*GetEqRequest)
		eq, err := svc.Get(ctx, gr.SubID)
		if err != nil {
			return nil, err
		}
		return responses.BaseResponse{
			Data: eq,
		}, nil
	}
}
