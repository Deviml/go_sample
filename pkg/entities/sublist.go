package entities

import "time"

type Sublist struct {
	ID                int               `json:"id"`
	Location          Location          `json:"location"`
	GeneralContractor GeneralContractor `json:"general_contractor"`
	ProjectName       string            `json:"complete_name"`
	PublishDate       time.Time         `json:"publish_date_format"`
	PublishDateUnix   int64             `json:"publish_date"`
	PurchasedAt       time.Time         `json:"-"`
}

type CompleteSublist struct {
	ID                int                      `json:"id"`
	Location          Location                 `json:"location"`
	GeneralContractor GeneralContractorSummary `json:"general_contractor"`
	ProjectName       string                   `json:"complete_name"`
	PublishDate       time.Time                `json:"publish_date_format"`
	PublishDateUnix   int64                    `json:"publish_date"`
	Companies         []SublistCompany         `json:"companies"`
}

type SublistCompany struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}
