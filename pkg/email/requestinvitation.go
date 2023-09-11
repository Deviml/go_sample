package email

import (
	"bytes"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
	"github.com/go-kit/kit/log"
	"html/template"
)

type RequestInvitation struct {
	logger log.Logger
	Name   string
	Phone  string
	Email  string
}

func NewRequestInvitation(logger log.Logger, name string, phone string, email string) *RequestInvitation {
	return &RequestInvitation{logger: logger, Name: name, Phone: phone, Email: email}
}

func (r RequestInvitation) Subject() string {
	return "General Contractor Invitation Request"
}

func (r RequestInvitation) Generate() string {
	cht := templates.RequestTemplate{
		Name:  r.Name,
		Phone: r.Phone,
		Email: r.Email,
	}
	t, _ := template.New("newInvitation").Parse(templates.RequestEmail)
	var tpl bytes.Buffer
	_ = t.Execute(&tpl, cht)
	return tpl.String()
}

func (r RequestInvitation) To() []string {
	return []string{"support@equiphunter.com"}
}
