package email

import (
	"bytes"
	"html/template"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
)

type ActionProposal struct {
	logger       log.Logger
	Destinations []string
	Proposal     models.Proposal
	webURL       string
}

func NewActionProposal(logger log.Logger, destinations []string, proposal models.Proposal, webURL string) *ActionProposal {
	return &ActionProposal{logger: logger, Destinations: destinations, Proposal: proposal, webURL: webURL}
}

func (c ActionProposal) Subject() string {
	q := c.Proposal
	if q.Status == models.RejectedStatus {
		return "Opps! Your proposal is rejected!!"
	} else {
		return "Your are the winner!!"
	}
}

func (c ActionProposal) Generate() string {
	//result := "Changes on your quotes!!<br>"
	var tpl bytes.Buffer
	q := c.Proposal
	if q.Status == models.RejectedStatus {
		cht := templates.RejectProposalTemplate{
			User:   q.WebUser.Username,
			EqName: q.Quote.EquipmentRequest.Equipment.Name,
		}

		t, _ := template.New("rejectProposal").Parse(templates.RejectProposal)
		_ = t.Execute(&tpl, cht)
	} else {
		cht := templates.AcceptProposalTemplate{
			User:       q.WebUser.Username,
			EqName:     q.Quote.EquipmentRequest.Equipment.Name,
			ProposalID: q.ProposalNumber,
		}

		t, _ := template.New("acceptProposal").Parse(templates.AcceptProposal)
		_ = t.Execute(&tpl, cht)

	}

	return tpl.String()
}

func (c ActionProposal) To() []string {
	return c.Destinations
}
