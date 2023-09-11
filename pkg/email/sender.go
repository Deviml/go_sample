package email

import (
	"context"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/ses"
	"github.com/go-kit/kit/log"
)

type Sender struct {
	logger       log.Logger
	sesClient    ses.Client
	baseFrontURL string
}

func NewSender(logger log.Logger, sesClient ses.Client, baseFrontURL string) *Sender {
	return &Sender{logger: logger, sesClient: sesClient, baseFrontURL: baseFrontURL}
}

func (s Sender) Send(ctx context.Context, email Email) error {
	return s.sesClient.SendEmail(ctx, email.To(), email.Generate(), email.Subject())
}

func (s Sender) SendNewQuote(ctx context.Context, to []string, name string, Category string, Amount int, City string, Zipcode string, SpecialRequest string, State string, FrontURL string, ID uint) error {
	return s.sesClient.SendEmailInBulkV2(ctx, to, name, Category, Amount, City, Zipcode, SpecialRequest, State, FrontURL, ID)
}

func (s Sender) SendForgotPassword(ctx context.Context, to string, newPassword string) error {
	return s.sesClient.SendForgotPassword(ctx, to, newPassword)
}

func (s Sender) SendConfirmation(ctx context.Context, to string, code string) error {
	e := NewConfirmation(to, code)
	return s.Send(ctx, e)
}

func (s Sender) SendWelcome(ctx context.Context, to string, name string) error {
	e := NewWelcome(to, name, s.baseFrontURL)
	return s.Send(ctx, e)
}

func (s Sender) SendActionProposal(ctx context.Context, to []string, q models.Proposal) error {
	e := NewActionProposal(s.logger, to, q, s.baseFrontURL)
	return s.Send(ctx, e)
}

func (s Sender) SendReceipt(ctx context.Context, to string, name string, purchases []templates.Purchase, total string) error {
	e := NewReceipt(to, name, purchases, total)
	return s.Send(ctx, e)
}

func (s Sender) SendNewSublist(ctx context.Context, to []string, ms []models.Sublist) error {
	e := NewNewSublist(s.logger, to, ms, s.baseFrontURL)
	return s.Send(ctx, e)
}

func (s Sender) SendSublistChange(ctx context.Context, to []string, ms models.Sublist, sc []models.SublistsCompany) error {
	e := NewChangeSublist(to, ms, sc, s.logger, s.baseFrontURL)
	return s.Send(ctx, e)
}

func (s Sender) SendInvitation(ctx context.Context, name string, phone string, email string) error {
	e := NewRequestInvitation(s.logger, name, phone, email)
	return s.Send(ctx, e)
}

func (s Sender) SendWelcomeWithCoupon(ctx context.Context, to string, name string) error {
	e := NewCouponWelcome(to, name, s.baseFrontURL)
	return s.Send(ctx, e)
}