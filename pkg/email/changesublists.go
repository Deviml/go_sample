package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
)

type Email interface {
	Subject() string
	Generate() string
	To() []string
}

type ChangeSublist struct {
	Destinations     []string
	Sublists         models.Sublist
	SublistCompanies []models.SublistsCompany
	logger           log.Logger
	webURL           string
}

func NewChangeSublist(destinations []string, sublists models.Sublist, sublistCompanies []models.SublistsCompany, logger log.Logger, webURL string) *ChangeSublist {
	return &ChangeSublist{Destinations: destinations, Sublists: sublists, SublistCompanies: sublistCompanies, logger: logger, webURL: webURL}
}

func (c ChangeSublist) Subject() string {
	return "Changes in your sublists!"
}

func (c ChangeSublist) To() []string {
	return c.Destinations
}

func (c ChangeSublist) Generate() string {
	cht := templates.ChangeSublistTemplate{
		Name:      c.Sublists.ProjectName,
		City:      c.Sublists.City.Name,
		State:     c.Sublists.City.State.Name,
		Zipcode:   c.Sublists.Zipcode,
		Companies: make([]templates.CompaniesTemplate, 0),
		Link:      fmt.Sprintf("%s/dashboard/sublists/%d", c.webURL, c.Sublists.ID),
	}
	for _, sc := range c.SublistCompanies {
		cht.Companies = append(cht.Companies, templates.CompaniesTemplate{
			Name:     sc.Company.CompanyName,
			Category: sc.CompanyCategory,
		})
	}

	var tpl bytes.Buffer
	t, _ := template.New("changeSub").Parse(templates.ChangeSublist)
	_ = t.Execute(&tpl, cht)
	return tpl.String()
}

func generateSublistEmail(s models.Sublist, sc []models.SublistsCompany) string {
	result := fmt.Sprintf("Project Name: %s<br>", s.ProjectName)
	for _, c := range sc {
		result += generateSublistCompany(c)
	}
	return result
}

func generateSublistCompany(sc models.SublistsCompany) string {
	result := fmt.Sprintf("Company: %s<br>", sc.Company.CompanyName)
	result += fmt.Sprintf("Category: %s<br>", sc.CompanyCategory)
	return result
}
