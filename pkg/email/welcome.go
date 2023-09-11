package email

import (
	"bytes"
	"html/template"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
)

type Welcome struct {
	Destination string
	Name        string
	URL         string
}

type CouponWelcome struct {
	Destination string
	Name        string
	URL         string
}

func NewCouponWelcome(destination string, name string, URL string) *CouponWelcome {
	return &CouponWelcome{Destination: destination, Name: name, URL: URL}
}

func NewWelcome(destination string, name string, URL string) *Welcome {
	return &Welcome{Destination: destination, Name: name, URL: URL}
}

func (w CouponWelcome) Subject() string {
	return "Welcome to Equiphunter"
}

func (w Welcome) Subject() string {
	return "Welcome to Equiphunter"
}

func (w CouponWelcome) Generate() string {
	cht := templates.WelcomeTemplate{
		Name: w.Name,
		URL:  w.URL,
	}
	t, _ := template.New("newWelcome").Parse(templates.WelcomeEmail)
	var tpl bytes.Buffer
	_ = t.Execute(&tpl, cht)
	return tpl.String()
}

func (w Welcome) Generate() string {
	cht := templates.WelcomeTemplate{
		Name: w.Name,
		URL:  w.URL,
	}
	t, _ := template.New("newWelcome").Parse(templates.WelcomeEmail)
	var tpl bytes.Buffer
	_ = t.Execute(&tpl, cht)
	return tpl.String()
}

func (w CouponWelcome) To() []string {
	return []string{w.Destination}
}

func (w Welcome) To() []string {
	return []string{w.Destination}
}
