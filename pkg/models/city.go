package models

import (
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"gorm.io/gorm"
	"strconv"
)

type City struct {
	gorm.Model
	Name      string
	StateID   uint
	State     State
	Latitude  float64
	Longitude float64
}

func (c City) ToEntity() entities.City {
	return entities.City{
		ID:   strconv.Itoa(int(c.ID)),
		Name: c.Name,
	}
}
