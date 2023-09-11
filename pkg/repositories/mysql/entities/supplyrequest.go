package entities

import "database/sql"

type SupplyRequest struct {
	Name           sql.NullString `json:"name"`
	Amount         sql.NullInt32  `json:"amount"`
	Supply         Supply         `json:"supply"`
	SpecialRequest sql.NullString
}
