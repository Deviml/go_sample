package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
)

type ChangeQuote struct {
	logger       log.Logger
	Destinations []string
	Quote        models.Quote
	webURL       string
}

func NewChangeQuote(logger log.Logger, destinations []string, quote models.Quote, webURL string) *ChangeQuote {
	return &ChangeQuote{logger: logger, Destinations: destinations, Quote: quote, webURL: webURL}
}

func (c ChangeQuote) Subject() string {
	return "Changes on your quotes!!"
}

func (c ChangeQuote) Generate() string {
	result := "Changes on your quotes!!<br>"
	q := c.Quote
	cht := templates.ChangeQuoteTemplate{
		City:    q.City.Name,
		State:   q.City.State.Name,
		Zipcode: q.Zipcode,
	}
	if q.SupplyRequest != nil {
		result += generateSupply(q)
		cht.Name = q.SupplyRequest.Supply.Name
		cht.Category = q.SupplyRequest.Supply.SupplyCategory.Name
		cht.SpecialRequest = q.SupplyRequest.SpecialRequest.String
		cht.Amount = q.SupplyRequest.Amount
		cht.Link = fmt.Sprintf("%s/dashboard/supplies/%d", c.webURL, q.ID)
	}

	if q.EquipmentRequest != nil {
		result += generateEquipment(q)
		cht.Name = q.EquipmentRequest.Equipment.Name
		cht.Category = q.EquipmentRequest.Equipment.EquipmentSubcategory.Name
		cht.SpecialRequest = q.EquipmentRequest.SpecialRequest.String
		cht.Amount = q.EquipmentRequest.Amount
		cht.Link = fmt.Sprintf("%s/dashboard/equipments/%d", c.webURL, q.ID)
	}

	t, _ := template.New("changeQuote").Parse(templates.ChangeQuote)
	var tpl bytes.Buffer
	_ = t.Execute(&tpl, cht)

	return tpl.String()
}

func (c ChangeQuote) To() []string {
	return c.Destinations
}
