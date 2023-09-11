package mysql

import (
	"context"
	"database/sql"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-kit/kit/log"
)

const (
	equipmentCategoriesTable = "EquipmentCategories"
)

type EquipmentCategoriesRepository struct {
	logger log.Logger
	client *sql.DB
}

func NewEquipmentCategoriesRepository(logger log.Logger, client *sql.DB) *EquipmentCategoriesRepository {
	return &EquipmentCategoriesRepository{logger: logger, client: client}
}

func (e EquipmentCategoriesRepository) List(ctx context.Context) ([]entities.EquipmentCategory, error) {
	query, args, err := e.makeListEquipmentCategoriesQuery(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := e.client.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()
	equipmentCategories := make([]entities.EquipmentCategory, 0)
	for rows.Next() {
		var equipmentCategory entities.EquipmentCategory
		scans := []interface{}{
			&equipmentCategory.ID,
			&equipmentCategory.Name,
		}
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		equipmentCategories = append(equipmentCategories, equipmentCategory)
	}
	return equipmentCategories, nil
}

func (e EquipmentCategoriesRepository) makeListEquipmentCategoriesQuery(ctx context.Context) (string, []interface{}, error) {
	return sq.Select("id", "name").From(equipmentCategoriesTable).Where("deleted_at IS NULL").Where("name != ''").ToSql()
}
