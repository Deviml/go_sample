package email

import (
	"bytes"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
	"html/template"
)

type Confirmation struct {
	Destination      string
	ConfirmationCode string
}

func NewConfirmation(destination string, confirmationCode string) *Confirmation {
	return &Confirmation{Destination: destination, ConfirmationCode: confirmationCode}
}

func (c Confirmation) Subject() string {
	return "Equiphunter Confirmation Code"
}

func (c Confirmation) Generate() string {
	cht := templates.ConfirmationTemplate{
		Code: c.ConfirmationCode,
	}
	t, _ := template.New("newConfirmation").Parse(templates.ConfirmationEmail)
	var tpl bytes.Buffer
	_ = t.Execute(&tpl, cht)
	return tpl.String()
}

func (c Confirmation) To() []string {
	return []string{c.Destination}
}
