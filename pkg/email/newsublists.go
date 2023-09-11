package email

import (
	"bytes"
	"fmt"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"html/template"
)

type NewSublist struct {
	logger       log.Logger
	Destinations []string
	Sublists     []models.Sublist
	link         string
}

func NewNewSublist(logger log.Logger, destinations []string, sublists []models.Sublist, link string) *NewSublist {
	return &NewSublist{logger: logger, Destinations: destinations, Sublists: sublists, link: link}
}

func (n NewSublist) Subject() string {
	return "There are new Sublists near you!!!"
}

func (n NewSublist) Generate() string {
	result := "These are the new sublists near your:<br>"
	for _, s := range n.Sublists {
		result += generateSummarySublist(s)
		cht := templates.NewSublistTemplate{
			ID:      int(s.ID),
			City:    s.City.Name,
			State:   s.City.State.Name,
			Zipcode: s.Zipcode,
			Link:    fmt.Sprintf("%s/sublists/?keywords=%d", n.link, s.ID),
		}
		t, _ := template.New("newSublist").Parse(templates.NewSublist)
		var tpl bytes.Buffer
		_ = t.Execute(&tpl, cht)
		return tpl.String()
	}
	return result
}

func generateSummarySublist(s models.Sublist) string {
	result := fmt.Sprintf("Project ID: %d<br>", s.ID)
	result += fmt.Sprintf("Location: %s<br>", s.Location.State)
	result += fmt.Sprintf("Zipcode: %d<br><br>", s.Location.Zipcode)
	return result
}

func (n NewSublist) To() []string {
	return n.Destinations
}
