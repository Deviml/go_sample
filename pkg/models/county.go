package models

import (
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"gorm.io/gorm"
	"strconv"
)

type State struct {
	gorm.Model
	Name string
}

func (s State) ToEntity() entities.State {
	return entities.State{
		ID:   strconv.Itoa(int(s.ID)),
		Name: s.Name,
	}
}

type County struct {
	gorm.Model
	Name    string
	StateID uint
	State   State
}

func (c County) ToEntity() entities.County {
	return entities.County{
		ID:   strconv.Itoa(int(c.ID)),
		Name: c.Name,
	}
}
