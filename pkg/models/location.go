package models

import (
	"strconv"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
)

type Location struct {
	ID      int
	Zipcode int
	City    string
	State   string
	Country string
}

func (l Location) ToEntity() entities.Location {
	return entities.Location{
		Zipcode: strconv.Itoa(l.Zipcode),
		City:    l.City,
		State:   l.State,
		Country: l.Country,
	}
}
