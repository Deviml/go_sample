package proposals

import (
	"context"
	"strconv"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/proposals"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type BuyerService struct {
	logger log.Logger
	db     *gorm.DB
	es     email.Sender
}

func NewBuyerService(logger log.Logger, db *gorm.DB, es email.Sender) *BuyerService {
	return &BuyerService{logger: logger, db: db, es: es}
}

func (e BuyerService) AcceptProposal(ctx context.Context, proposalID string) error {
	var proposal models.Proposal
	var quote models.Quote
	result := e.db.
		Find(&proposal, "id = ?", proposalID)
	if result.Error != nil {
		e.logger.Log("resultError", result.Error)
		return result.Error
	}
	e.db.Model(&proposal).Update("status", models.ApprovedStatus)
	result = e.db.Model(&quote).Where("id = ?", &proposal.QuoteID).Update("status", models.ServedStatus)
	if result.Error != nil {
		e.logger.Log("err", result.Error.Error())
		return result.Error
	}
	e.db.Preload("WebUser").
		Preload("Quote.EquipmentRequest.Equipment").
		Preload("Quote.WebUser").
		Find(&proposal, "id = ?", proposalID)
	err := e.notifyChange(ctx, proposal)
	if err != nil {
		e.logger.Log("err", err.Error())
		return err
	}

	//reject all other APIs
	var proposals []models.Proposal
	result = e.db.Where("proposals.status = ?", models.NewStatus).
		Find(&proposals, "proposals.quote_id = ?", proposal.QuoteID)
	if result != nil && result.Error != nil {
		return nil
	}
	for _, mq := range proposals {
		e.RejectProposal(ctx, strconv.Itoa(int(mq.ID)))
	}

	return nil
}

func (e BuyerService) notifyChange(ctx context.Context, proposal models.Proposal) error {
	user := proposal.WebUser
	emails := make([]string, 0)
	emails = append(emails, user.Username)
	return e.es.SendActionProposal(ctx, emails, proposal)
}

func (e BuyerService) RejectProposal(ctx context.Context, proposalID string) error {
	var proposal models.Proposal
	result := e.db.
		Find(&proposal, "id = ?", proposalID)
	if result.Error != nil {
		return result.Error
	}

	e.db.Model(&proposal).Where("id = ?", proposalID).Update("status", models.RejectedStatus)
	e.db.Preload("WebUser").Preload("Quote.EquipmentRequest.Equipment").Find(&proposal, "id = ?", proposalID)
	err := e.notifyChange(ctx, proposal)
	if err != nil {
		e.logger.Log("err", err.Error())
		return err
	}
	return nil
}

func (u BuyerService) Get(ctx context.Context, uqRequest proposals.BuyerProposalRequest) ([]entities.Proposal, error) {
	var user models.WebUser
	result := u.db.Find(&user, "id = ?", uqRequest.UserID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}

	var proposals []models.Proposal
	result = u.db.
		Preload("WebUser").
		Preload("Quote.EquipmentRequest.Equipment").
		Find(&proposals, "Quotes.web_user_id = ?", uqRequest.UserID).Where("proposals.status = ?", uqRequest.Status)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	eQ := make([]entities.Proposal, 0)
	for _, mq := range proposals {
		eQ = append(eQ, mq.ToCompleteProposalEntity())
	}
	return eQ, nil
}
