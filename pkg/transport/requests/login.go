package requests

import (
	"context"
	"encoding/json"
	http2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/errors/http"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func DecodeLoginRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var loginRequest LoginRequest
	err = json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		return nil, http2.BadRequest{}
	}
	return &loginRequest, nil
}
