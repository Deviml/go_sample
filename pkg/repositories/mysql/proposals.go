package mysql

import (
	"context"
	"database/sql"
	"fmt"

	domainEntities "github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/repositories/mysql/entities"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-kit/kit/log"
)

const (
	ProposalTable = "proposals"

//	dateTimeFormat = "2006-01-02 15:04:05"
)

type ProposalsRepositories struct {
	logger log.Logger
	client *sql.DB
}

func NewProposalsRepositories(logger log.Logger, client *sql.DB) *ProposalsRepositories {
	return &ProposalsRepositories{logger: logger, client: client}
}

func (q ProposalsRepositories) List(ctx context.Context, userID string, paginationQuery domainEntities.PaginationQuery, userType string) ([]entities.Proposal, error) {
	query, args, err := q.makeListProposalsQuery(ctx, userID, paginationQuery, userType)
	q.logger.Log("query", query)
	q.logger.Log("args", fmt.Sprintf("%v", args))
	q.logger.Log("userid", userID)
	if err != nil {
		q.logger.Log("error", err)
		return nil, err
	}

	rows, err := q.client.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	proposals := make([]entities.Proposal, 0)
	for rows.Next() {
		var proposal entities.Proposal
		scans := []interface{}{
			&proposal.ID,
			&proposal.Make,
			&proposal.Model,
			&proposal.Year,
			&proposal.Vin,
			&proposal.Serial,
			&proposal.EqHours,
			&proposal.Condition,
			&proposal.SalePrice,
			&proposal.AvailableDate,
			&proposal.Description,
			&proposal.Comments,
			&proposal.Specifications,
			&proposal.Videos,
			&proposal.Pics,
			&proposal.Status,
			&proposal.Freight,
			&proposal.Tax,
			&proposal.Fees,
			&proposal.EquipmentName,
			&proposal.ProposalNumber,
		}
		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		proposals = append(proposals, proposal)
	}
	return proposals, nil
}

func (q ProposalsRepositories) makeListProposalsQuery(ctx context.Context, userID string, paginationQuery domainEntities.PaginationQuery, userType string) (string, []interface{}, error) {
	query := sq.Select(
		"DISTINCT(proposals.id)",
		"proposals.make",
		"proposals.eq_model",
		"proposals.year",
		"proposals.serial",
		"proposals.vin",
		"proposals.eq_hours",
		"proposals.condition",
		"proposals.sale_price",
		"proposals.available_date",
		"proposals.description",
		"proposals.comments",
		"proposals.specifications",
		"proposals.videos",
		"proposals.pics",
		"proposals.status",
		"proposals.freight",
		"proposals.tax",
		"proposals.fees",
		"equipments.name",
		"proposals.proposal_number",
	).
		From(ProposalTable).
		LeftJoin("Quotes ON Quotes.id = proposals.quote_id").
		LeftJoin("EquipmentRequests ON Quotes.equipment_request_id = EquipmentRequests.id").
		LeftJoin("equipments ON EquipmentRequests.equipment_id = equipments.id")
	query = applyListPagination(query, paginationQuery)
	query = q.applyFiltersListProposals(query, userID, userType)
	query = q.applySortListProposals(query)
	return query.ToSql()
}

func (q ProposalsRepositories) applySortListProposals(query sq.SelectBuilder) sq.SelectBuilder {
	query = query.OrderBy("proposals.quote_id desc", "proposals.created_at desc")
	return query
}

func (q ProposalsRepositories) applyFiltersListProposals(query sq.SelectBuilder, userID string, userType string) sq.SelectBuilder {
	query = query.Where(sq.Eq{"proposals.deleted_at": nil}).Where(sq.Eq{"Quotes.deleted_at": nil})
	q.logger.Log(userType)
	if userID != "" {
		if userType == "buyer" {
			query = query.Where(sq.Eq{"Quotes.web_user_id": userID})
		} else {
			query = query.Where(sq.Eq{"proposals.web_user_id": userID})
		}
	}
	return query
}

func (q ProposalsRepositories) makeFetchCountProposalsQuery(ctx context.Context, userID string, userType string) (string, []interface{}, error) {
	query := sq.Select("COUNT(DISTINCT(proposals.id))").
		From(ProposalTable).
		LeftJoin("Quotes ON Quotes.id = proposals.quote_id")
	query = q.applyFiltersListProposals(query, userID, userType)
	return query.ToSql()
}

func (q ProposalsRepositories) FetchCount(ctx context.Context, userID string, userType string) (int, error) {
	query, args, err := q.makeFetchCountProposalsQuery(ctx, userID, userType)
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
