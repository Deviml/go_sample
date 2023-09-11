package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/go-kit/kit/log"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
)

type NewQuotes struct {
	logger       log.Logger
	Quotes       []models.Quote
	Destinations []string
	Link         string
}

func NewNewQuotes(logger log.Logger, quotes []models.Quote, destinations []string, link string) *NewQuotes {
	return &NewQuotes{logger: logger, Quotes: quotes, Destinations: destinations, Link: link}
}

func (n NewQuotes) Subject() string {
	return "New Quotes Around Your Area"
}

func (n NewQuotes) To() []string {
	return n.Destinations
}

func (n NewQuotes) Generate() string {
	result := ""
	for _, q := range n.Quotes {
		cht := templates.NewQuoteTemplate{
			ID:      q.ID,
			City:    q.City.Name,
			State:   q.City.State.Name,
			Zipcode: q.Zipcode,
			Link:    fmt.Sprintf("%s/quotes/?keywords=%d", n.Link, q.ID),
		}
		if q.SupplyRequest != nil {
			cht.Name = q.SupplyRequest.Supply.Name
			cht.Category = q.SupplyRequest.Supply.SupplyCategory.Name
			cht.SpecialRequest = q.SupplyRequest.SpecialRequest.String
			cht.Amount = q.SupplyRequest.Amount
		}

		if q.EquipmentRequest != nil {
			cht.Name = q.EquipmentRequest.Equipment.Name
			cht.Category = q.EquipmentRequest.Equipment.EquipmentSubcategory.Name
			cht.SpecialRequest = q.EquipmentRequest.SpecialRequest.String
			cht.Amount = q.EquipmentRequest.Amount
		}
		var tpl bytes.Buffer
		t, _ := template.New("newQuote").Parse(templates.NewQuote)
		_ = t.Execute(&tpl, cht)
		return tpl.String()
	}
	return result
}

func generateSupply(q models.Quote) string {
	result := fmt.Sprintf("Supply: %s<br>", q.SupplyRequest.Supply.Name)
	result += fmt.Sprintf("Supply Category: %s<br>", q.SupplyRequest.Supply.SupplyCategory.Name)
	result += fmt.Sprintf("Amount: %d<br><br>", q.SupplyRequest.Amount)
	result += fmt.Sprintf("Special Request: %s<br><br>", q.SupplyRequest.SpecialRequest.String)
	return result
}

func generateEquipment(q models.Quote) string {
	result := fmt.Sprintf("Equipment: %s<br>", q.EquipmentRequest.Equipment.Name)
	result += fmt.Sprintf("Equipment Category: %s<br>", q.EquipmentRequest.Equipment.EquipmentSubcategory.EquipmentCategory.Name)
	result += fmt.Sprintf("Equipment Sub Category: %s<br><br>", q.EquipmentRequest.Equipment.EquipmentSubcategory.Name)
	result += fmt.Sprintf("Special Request: %s<br><br>", q.EquipmentRequest.SpecialRequest.String)
	return result
}
