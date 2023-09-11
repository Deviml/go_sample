package requests

import (
	"context"
	"net/http"
)

type ListCompaniesRequest struct {
	Keyword string
}

func DecodeListCompaniesRequest(ctx context.Context, r *http.Request) (response interface{}, err error) {
	return &ListCompaniesRequest{
		Keyword: r.URL.Query().Get("keywords"),
	}, nil
}
