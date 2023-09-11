package responses

import "net/http"

type NoContent struct {
}

func (n NoContent) StatusCode() int {
	return http.StatusNoContent
}
