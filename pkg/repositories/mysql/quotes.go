package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	domainEntities "github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/filters"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities/sorts"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/repositories/mysql/entities"
	sq "github.com/Masterminds/squirrel"

	"github.com/go-kit/kit/log"
)

const (
	quotesTable    = "Quotes"
	dateTimeFormat = "2006-01-02 15:04:05"
)

type QuotesRepositories struct {
	logger log.Logger
	client *sql.DB
}

func NewQuotesRepositories(logger log.Logger, client *sql.DB) *QuotesRepositories {
	return &QuotesRepositories{logger: logger, client: client}
}

func (q QuotesRepositories) List(ctx context.Context, filter filters.ListQuotes, sort sorts.ListQuotes, paginationQuery domainEntities.PaginationQuery) ([]entities.Quote, error) {
	query, args, err := q.makeListQuotesQuery(ctx, filter, sort, paginationQuery)
	q.logger.Log("query", query)
	q.logger.Log("args", fmt.Sprintf("%v", args))
	if err != nil {
		q.logger.Log("error", err)
		return nil, err
	}

	rows, err := q.client.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	quotes := make([]entities.Quote, 0)
	for rows.Next() {
		var quote entities.Quote
		scans := []interface{}{
			&quote.ID,
			&quote.Location.Zipcode,
			&quote.Location.City,
			&quote.Location.Country,
			&quote.Location.State,
			&quote.EquipmentRequest.Equipment.Name,
			&quote.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory.Name,
			&quote.EquipmentRequest.Equipment.EquipmentSubcategory.Name,
			&quote.EquipmentRequest.Amount,
			&quote.EquipmentRequest.SpecialRequest,
			&quote.EquipmentRequest.ContractPreference,
			&quote.EquipmentRequest.EndDateFormatted,
			&quote.EquipmentRequest.StartDateFormatted,
			&quote.EquipmentRequest.ExpirationDate,
			&quote.SupplyRequest.Supply.Name,
			&quote.SupplyRequest.Amount,
			&quote.SupplyRequest.Supply.SupplyCategory.Name,
			&quote.ProductName,
			&quote.Status,
		}
		if filter.UserID != "" {
			scans = append(scans, &quote.PurchasedAt)
		}
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}
	return quotes, nil
}

func (q QuotesRepositories) makeListQuotesQuery(ctx context.Context, filter filters.ListQuotes, sort sorts.ListQuotes, paginationQuery domainEntities.PaginationQuery) (string, []interface{}, error) {
	query := sq.Select(
		"DISTINCT(Quotes.id)",
		"Quotes.zipcode",
		"cities.name",
		"locations.country",
		"states.name",
		"equipments.name",
		"EquipmentCategories.name",
		"EquipmentSubcategories.name",
		"EquipmentRequests.amount",
		"EquipmentRequests.special_request",
		"EquipmentRequests.contract_preference",
		"EquipmentRequests.rent_to",
		"EquipmentRequests.rent_from",
		"EquipmentRequests.expiration_date",
		"supplies.name",
		"SupplyRequests.amount",
		"SupplyCategories.name",
		"COALESCE(EquipmentCategories.name, SupplyCategories.name) as productName",
		"Quotes.status",
	).
		From(quotesTable).
		LeftJoin("EquipmentRequests ON EquipmentRequests.id = Quotes.equipment_request_id").
		LeftJoin("SupplyRequests ON SupplyRequests.id = Quotes.supply_request_id").
		LeftJoin("equipments ON EquipmentRequests.equipment_id = equipments.id").
		LeftJoin("supplies ON SupplyRequests.supply_id = supplies.id").
		LeftJoin("EquipmentSubcategories ON equipments.equipment_subcategory_id = EquipmentSubcategories.id").
		LeftJoin("EquipmentCategories ON equipments.equipment_category_id = EquipmentCategories.id").
		LeftJoin("SupplyCategories ON SupplyCategories.id = supplies.supply_category_id").
		Join("locations ON EquipmentRequests.location_id = locations.id OR SupplyRequests.location_id = locations.id").
		Join("cities ON cities.id = Quotes.city_id").
		Join("states ON cities.state_id = states.id").
		GroupBy("Quotes.id").
		GroupBy("EquipmentSubcategories.name")
	query = applyListPagination(query, paginationQuery)
	query = q.applyFiltersListQuotes(query, filter)
	query = q.applySortListQuotes(query, sort)
	a, b, c := query.ToSql()
	return a, b, c
}

func (q QuotesRepositories) applySortListQuotes(query sq.SelectBuilder, sort sorts.ListQuotes) sq.SelectBuilder {
	switch sort.Criteria {
	case "date":
		query = query.OrderBy("Quotes.created_at DESC")
	case "name":
		query = query.OrderBy("Quotes.created_at DESC")
	case "location":
		query = query.OrderBy("Quotes.created_at DESC")
	default:
		query = query.OrderBy("Quotes.created_at DESC")
	}
	return query
}

func (q QuotesRepositories) applyFiltersListQuotes(query sq.SelectBuilder, filter filters.ListQuotes) sq.SelectBuilder {
	query = query.Where(sq.Eq{"Quotes.deleted_at": nil}).Where(sq.Eq{"EquipmentCategories.deleted_at": nil}).Where(sq.Eq{"SupplyCategories.deleted_at": nil})

	if filter.StateID != 0 && filter.CityID == 0 {
		query = query.Where(sq.Eq{"cities.state_id": filter.StateID})
	}

	if filter.CityID != 0 {
		query = query.Where(sq.Eq{"Quotes.city_id": filter.CityID})
	}

	if filter.Zipcode != "" {
		query = query.Where(sq.Like{"Quotes.zipcode": fmt.Sprintf("%%%s%%", filter.Zipcode)})
	}

	if filter.UserID != "" {
		query = query.Join("payment_details ON payment_details.quote_id = Quotes.id").
			Join("payments ON payments.id = payment_details.payment_id").
			Join("web_users ON web_users.id = payments.web_user_id").
			Columns("payments.created_at").
			Where(sq.Eq{"web_users.id": filter.UserID})
	}

	if filter.OnlyEquipment {
		query = query.Where(sq.NotEq{"Quotes.equipment_request_id": nil})
	} else if filter.OnlySupply {
		query = query.Where(sq.NotEq{"Quotes.supply_request_id": nil})
	}

	if filter.NotExpired && filter.NotServed {
		query = query.Where(sq.LtOrEq{"Quotes.status": 1})
	} else if filter.NotServed {
		query = query.Where(sq.NotEq{"Quotes.status": 2})
	} else if filter.NotExpired {
		query = query.Where(sq.NotEq{"Quotes.status": 3})
	}

	if len(filter.Keywords) > 0 {
		or := sq.Or{}
		
		for _, keyword := range filter.Keywords {
			// re := regexp.MustCompile(`[,|.|;|:|/|\\]`)
			// re := regexp.MustCompile(`[,]`)
			// keyword = re.ReplaceAllString(keyword, "")
			// keyword = strings.Replace(keyword, `,`, "", -1)
			or = append(or, sq.Like{"equipments.name": fmt.Sprintf("%%%s%%", keyword)})
			or = append(or, sq.Like{"supplies.name": fmt.Sprintf("%%%s%%", keyword)})
			or = append(or, sq.Eq{"Quotes.id": keyword})
			or = append(or, sq.Like{"cities.name": fmt.Sprintf("%%%s%%", keyword)})
			or = append(or, sq.Like{"Quotes.zipcode": fmt.Sprintf("%%%s%%", keyword)})
			or = append(or, sq.Like{"states.name": fmt.Sprintf("%%%s%%", keyword)})
		}
		if len(or) > 0 {
			query = query.Where(or)
		}
	}

	if len(filter.ContractPreferences) > 0 {
		or := sq.Or{}
		for _, contractPreference := range filter.ContractPreferences {
			or = append(or, sq.Eq{"EquipmentRequests.contract_preference": contractPreference})
		}
		if len(or) > 0 {
			query = query.Where(or)
		}
	}

	if len(filter.SupplyCategories) > 0 {
		or := sq.Or{}
		for _, categoryID := range filter.SupplyCategories {
			if categoryID == "" {
				continue
			}
			or = append(or, sq.Eq{"SupplyCategories.id": categoryID})
		}
		if len(or) > 0 {
			query = query.Where(or)
		}
	}

	if len(filter.EquipmentCategories) > 0 {
		or := sq.Or{}
		for _, categoryID := range filter.EquipmentCategories {
			if categoryID == "" {
				continue
			}
			or = append(or, sq.Eq{"EquipmentCategories.id": categoryID})
		}
		if len(or) > 0 {
			query = query.Where(or)
		}
	}

	if filter.EquipmentSubcategory != "" {
		query = query.Where(sq.Eq{"EquipmentSubcategories.id": filter.EquipmentSubcategory})
	}

	if filter.RentTo != 0 {
		query = query.Where(sq.LtOrEq{"EquipmentRequests.rent_to": time.Unix(filter.RentTo, 0).Format(dateTimeFormat)})
	}

	if filter.RentFrom != 0 {
		query = query.Where(sq.GtOrEq{"EquipmentRequests.rent_from": time.Unix(filter.RentFrom, 0).Format(dateTimeFormat)})
	}
	return query
}

func (q QuotesRepositories) makeFetchCountQuotesQuery(ctx context.Context, filter filters.ListQuotes) (string, []interface{}, error) {
	query := sq.Select("COUNT(DISTINCT(Quotes.id))").
		From(quotesTable).
		LeftJoin("EquipmentRequests ON EquipmentRequests.id = Quotes.equipment_request_id").
		LeftJoin("SupplyRequests ON SupplyRequests.id = Quotes.supply_request_id").
		LeftJoin("equipments ON EquipmentRequests.equipment_id = equipments.id").
		LeftJoin("supplies ON SupplyRequests.supply_id = supplies.id").
		LeftJoin("EquipmentSubcategories ON equipments.equipment_subcategory_id = EquipmentSubcategories.id").
		LeftJoin("EquipmentCategories ON equipments.equipment_category_id = EquipmentCategories.id").
		LeftJoin("SupplyCategories ON SupplyCategories.id = supplies.supply_category_id").
		Join("locations ON EquipmentRequests.location_id = locations.id OR SupplyRequests.location_id = locations.id").
		Join("cities ON cities.id = Quotes.city_id").
		Join("states on cities.state_id = states.id")
	query = q.applyFiltersListQuotes(query, filter)
	return query.ToSql()
}

func (q QuotesRepositories) FetchCount(ctx context.Context, filter filters.ListQuotes) (int, error) {
	query, args, err := q.makeFetchCountQuotesQuery(ctx, filter)
	q.logger.Log("query", query)
	q.logger.Log("args", fmt.Sprintf("%v", args))

	if err != nil {
		return 0, err
	}
	rows, err := q.client.Query(query, args...)
	if err != nil {
		return 0, err
	}

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
