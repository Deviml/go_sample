package http

import (
	"net/http"
)

type BadRequest struct{}

func (b BadRequest) Error() string {
	return "bad request"
}

func (b BadRequest) StatusCode() int {
	return http.StatusBadRequest
}
