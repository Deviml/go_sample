package email

import (
	"bytes"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email/templates"
	"html/template"
)

type Receipt struct {
	Destination string
	Name        string
	Purchases   []templates.Purchase
	Total       string
}

func NewReceipt(destination string, name string, purchases []templates.Purchase, total string) *Receipt {
	return &Receipt{Destination: destination, Name: name, Purchases: purchases, Total: total}
}

func (r Receipt) Subject() string {
	return "Equiphunter Receipt"
}

func (r Receipt) Generate() string {
	cht := templates.ReceiptTemplate{
		Name:      r.Name,
		Purchases: r.Purchases,
		Total:     r.Total,
	}
	t, _ := template.New("newReceipt").Parse(templates.ReceiptEmail)
	var tpl bytes.Buffer
	_ = t.Execute(&tpl, cht)
	return tpl.String()
}

func (r Receipt) To() []string {
	return []string{r.Destination}
}
