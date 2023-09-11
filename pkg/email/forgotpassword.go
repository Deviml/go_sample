package email

import (
	"bytes"
	"html/template"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
)

type ForgotPassword struct {
	Destination string
	Password    string
}

func NewForgotPassword(destination string, password string) *ForgotPassword {
	return &ForgotPassword{Destination: destination, Password: password}
}

func (f ForgotPassword) Subject() string {
	return "Equiphunter Password Reset"
}

func (f ForgotPassword) Generate() string {
	cht := templates.ForgotPasswordTemplate{
		Password: f.Password,
	}
	t, _ := template.New("newPassword").Parse(templates.ForgotPasswordEmail)
	var tpl bytes.Buffer
	_ = t.Execute(&tpl, cht)
	return tpl.String()
}

func (f ForgotPassword) To() []string {
	return []string{f.Destination}
}
