package mysql

import (
	"context"
	"database/sql"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-kit/kit/log"
)

const (
	equipmentSubcategoryTable = "EquipmentSubcategories"
)

type EquipmentSubcategoriesRepository struct {
	logger log.Logger
	client *sql.DB
}

func NewEquipmentSubcategoriesRepository(logger log.Logger, client *sql.DB) *EquipmentSubcategoriesRepository {
	return &EquipmentSubcategoriesRepository{logger: logger, client: client}
}

func (e EquipmentSubcategoriesRepository) GetList(ctx context.Context, equipmentCategoryID string) ([]entities.EquipmentSubcategory, error) {
	query, args, err := e.makeGetQuery(equipmentCategoryID)
	if err != nil {
		return nil, err
	}
	rows, err := e.client.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subcategories []entities.EquipmentSubcategory
	for rows.Next() {
		var subcategory entities.EquipmentSubcategory
		scans := []interface{}{
			&subcategory.ID,
			&subcategory.Name,
		}
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		subcategories = append(subcategories, subcategory)
	}
	return subcategories, nil
}

func (e EquipmentSubcategoriesRepository) makeGetQuery(equipmentCategoryID string) (string, []interface{}, error) {
	query := sq.Select("id", "name").From(equipmentSubcategoryTable)
	if equipmentCategoryID != "" {
		query = query.Where(sq.Eq{"equipment_category_id": equipmentCategoryID})
	}
	return query.ToSql()
}
