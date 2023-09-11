package entities

import (
	"database/sql"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
)

type Supply struct {
	ID             string         `json:"id"`
	Name           sql.NullString `json:"name"`
	SupplyCategory SupplyCategory `json:"supply_category"`
}

func (s Supply) ToEntity() entities.Supply {
	return entities.Supply{
		ID:             s.ID,
		Name:           s.Name.String,
		SupplyCategory: s.SupplyCategory.ToEntity(),
	}
}
