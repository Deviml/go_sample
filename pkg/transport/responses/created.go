package responses

import "net/http"

type Created struct {
}

func (c Created) StatusCode() int {
	return http.StatusCreated
}
