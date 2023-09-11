package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/filters"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/sorts"
	sq "github.com/Masterminds/squirrel"
	sq2 "github.com/huandu/go-sqlbuilder"

	"github.com/go-kit/kit/log"
)

const (
	sublistsTable = "Sublists"
)

type SublistsRepository struct {
	logger log.Logger
	db     *sql.DB
}

func NewSublistsRepository(logger log.Logger, db *sql.DB) *SublistsRepository {
	return &SublistsRepository{logger: logger, db: db}
}

func (s SublistsRepository) FetchSublists(ctx context.Context, filters filters.ListSublists, sorts sorts.ListSublists, pagination entities.PaginationQuery) ([]entities.Sublist, error) {
	query, args, err := makeListQuery(filters, sorts, pagination)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	sublists := make([]entities.Sublist, 0)
	for rows.Next() {
		var sublist entities.Sublist
		scans := []interface{}{
			&sublist.ID,
			&sublist.ProjectName,
			&sublist.PublishDate,
			&sublist.GeneralContractor.Name,
			&sublist.Location.Zipcode,
			&sublist.Location.City,
			&sublist.Location.State,
			&sublist.Location.Country,
			&sublist.GeneralContractor.Company.Name,
		}
		if filters.UserID != "" {
			scans = append(scans, &sublist.PurchasedAt)
		}
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		sublists = append(sublists, sublist)
	}
	return sublists, nil
}

func makeListQuery(filters filters.ListSublists, sorts sorts.ListSublists, pagination entities.PaginationQuery) (string, []interface{}, error) {
	query := sq.Select(
		"Sublists.id",
		"Sublists.Project_Name",
		"Sublists.created_at",
		"web_users.full_name",
		"Sublists.zipcode",
		"cities.name",
		"states.name",
		"locations.country",
		"companies.company_name",
	).From(sublistsTable)

	query = applyListJoins(query)
	query = applyListFilters(query, filters)
	query = applyListSorts(query, sorts)
	query = applyListPagination(query, pagination)
	return query.ToSql()
}

func applyListJoins(query sq.SelectBuilder) sq.SelectBuilder {
	return query.
		Join("web_users on web_users.id = Sublists.web_user_id").
		Join("general_contractors on general_contractors.id = web_users.profile_id").
		Join("locations ON  locations.id = Sublists.location_id").
		Join("companies ON general_contractors.company_id = companies.id").
		Join("cities ON Sublists.city_id = cities.id").
		Join("states ON cities.state_id = states.id")
}

func applyListPagination(query sq.SelectBuilder, pagination entities.PaginationQuery) sq.SelectBuilder {
	if pagination.PerPage != 0 {
		perPage := uint64(pagination.PerPage)
		query = query.Limit(perPage)
		if pagination.Page != 0 {
			page := uint64(pagination.Page)
			query = query.Offset((page - 1) * perPage)
		}
	}
	return query
}

func applyNewListPagination(query *sq2.SelectBuilder, pagination entities.PaginationQuery) *sq2.SelectBuilder {
	if pagination.PerPage != 0 {
		perPage := int(pagination.PerPage)
		query = query.Limit(perPage)
		if pagination.Page != 0 {
			page := int(pagination.Page)
			query = query.Offset((page - 1) * perPage)
		}
	}
	return query
}

func applyListSorts(query sq.SelectBuilder, sorts sorts.ListSublists) sq.SelectBuilder {
	switch sorts.Criteria {
	case "date":
		query = query.OrderBy("Sublists.created_at DESC")
	case "distance":
		distanceSort := fmt.Sprintf("3956 * 2 * ASIN(SQRT(\nPOWER(SIN((%f - abs(cities.latitude)) * pi()/180 / 2), 2) + COS(%f * pi()/180 ) * COS(abs(cities.latitude) * pi()/180) * POWER(SIN((%f - cities.longitude) *\npi()/180 / 2), 2) ))", sorts.Latitude, sorts.Latitude, sorts.Longitude)
		query = query.OrderBy(distanceSort)
	}
	return query
}

func applyListFilters(query sq.SelectBuilder, filters filters.ListSublists) sq.SelectBuilder {
	query = query.Where(sq.Eq{"Sublists.deleted_at": nil})

	if filters.StateID != 0 && filters.CityID == 0 {
		query = query.Where(sq.Eq{"cities.state_id": filters.CityID})
	}

	if filters.CityID != 0 {
		query = query.Where(sq.Eq{"cities.id": filters.CityID})
	}

	if filters.Zipcode != "" {
		query = query.Where(sq.Like{"Sublists.zipcode": fmt.Sprintf("%%%s%%", filters.Zipcode)})
	}

	if len(filters.Keywords) > 0 {
		or := sq.Or{}
		for _, keyword := range filters.Keywords {
			or = append(or, sq.Eq{"Sublists.id": keyword})
			or = append(or, sq.Like{"Sublists.Project_Name": fmt.Sprintf("%%%s%%", keyword)})
			or = append(or, sq.Like{"companies.company_Name": fmt.Sprintf("%%%s%%", keyword)})
		}
		if len(or) > 0 {
			query = query.Where(or)
		}
	}

	if filters.From != "" {
		query = query.Where(sq.GtOrEq{"Sublists.created_at": filters.From})
	}

	if filters.To != "" {
		query = query.Where(sq.LtOrEq{"Sublists.created_at": filters.To})
	}

	if filters.UserID != "" {
		query = query.
			Columns("payment_details.created_at").
			Join("web_user_sublists ON web_user_sublists.sublist_id = Sublists.id").
			LeftJoin("payment_details ON payment_details.sublist_id = Sublists.id").
			GroupBy("Sublists.id").Where(sq.Eq{"web_user_sublists.web_user_id": filters.UserID})
	}
	return query
}
func (s SublistsRepository) FetchCount(ctx context.Context, filters filters.ListSublists, sorts sorts.ListSublists) (int, error) {
	query, args, err := makeListCountQuery(filters)
	if err != nil {
		return 0, err
	}

	rows, err := s.db.Query(query, args...)
	defer func() {
		_ = rows.Close()
	}()
	var total int
	for rows.Next() {
		err = rows.Scan(&total)
		if err != nil {
			return 0, err
		}
		break
	}
	return total, nil
}

func makeListCountQuery(filters filters.ListSublists) (string, []interface{}, error) {
	query := sq.Select("COUNT(Sublists.id)").From(sublistsTable)
	query = applyListJoins(query)
	query = applyListFilters(query, filters)
	return query.ToSql()
}
