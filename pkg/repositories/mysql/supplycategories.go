package mysql

import (
	"context"
	"database/sql"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-kit/kit/log"
)

const (
	supplyCategoriesTable = "SupplyCategories"
)

type SupplyCategoriesRepository struct {
	logger log.Logger
	client *sql.DB
}

func NewSupplyCategoriesRepository(logger log.Logger, client *sql.DB) *SupplyCategoriesRepository {
	return &SupplyCategoriesRepository{logger: logger, client: client}
}

func (s SupplyCategoriesRepository) List(ctx context.Context) ([]entities.SupplyCategory, error) {
	query, args, err := s.makeListSupplyCategoriesQuery(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := s.client.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()
	supplyCategories := make([]entities.SupplyCategory, 0)
	for rows.Next() {
		var supplyCategory entities.SupplyCategory
		scans := []interface{}{
			&supplyCategory.ID,
			&supplyCategory.Name,
		}
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		supplyCategories = append(supplyCategories, supplyCategory)
	}
	return supplyCategories, nil
}

func (s SupplyCategoriesRepository) makeListSupplyCategoriesQuery(ctx context.Context) (string, []interface{}, error) {
	return sq.Select("id", "name").From(supplyCategoriesTable).Where("deleted_at IS NULL").Where("name != ''").ToSql()
}
