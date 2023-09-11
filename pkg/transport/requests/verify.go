package requests

import (
	"context"
	"encoding/json"
	http2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/errors/http"
	"net/http"
)

type VerifyRequest struct {
	Code string `json:"code"`
}

func DecodeVerifyRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	var verify VerifyRequest
	err = json.NewDecoder(r.Body).Decode(&verify)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &verify, nil
}
