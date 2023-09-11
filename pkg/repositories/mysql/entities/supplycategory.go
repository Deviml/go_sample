package entities

import (
	"database/sql"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
)

type SupplyCategory struct {
	ID   string         `json:"id"`
	Name sql.NullString `json:"name"`
}

func (sc SupplyCategory) ToEntity() entities.SupplyCategory {
	return entities.SupplyCategory{
		ID:   sc.ID,
		Name: sc.Name.String,
	}
}
