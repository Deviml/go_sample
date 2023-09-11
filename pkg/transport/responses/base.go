package responses

import (
	"net/http"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
)

type CollectionResponse struct {
	Data       interface{}          `json:"data"`
	Pagination *entities.Pagination `json:"pagination"`
}

func (b CollectionResponse) StatusCode() int {
	return http.StatusOK
}

type BaseResponse struct {
	Data interface{} `json:"data"`
}

func (b BaseResponse) StatusCode() int {
	return http.StatusOK
}

type SuccessMessage struct {
	Message string `json:"message"`
}

func (b SuccessMessage) StatusCode() int {
	return http.StatusOK
}
