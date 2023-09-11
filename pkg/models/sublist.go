package models

import (
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"gorm.io/gorm"
	"strconv"
)

type Sublist struct {
	gorm.Model
	ProjectName       string `gorm:"column:Project_Name"`
	LocationID        int
	Location          Location
	WebUserID         int
	WebUser           WebUser
	Companies         []Company `gorm:"many2many:sublists_companies;"`
	SublistsCompanies []SublistsCompany
	WebUsers          []WebUser `gorm:"many2many:web_user_sublists;"`
	Zipcode           string
	CityID            int
	City              *City
	CountyID          *int
	County            *County
}

type SublistsCompany struct {
	gorm.Model
	SublistID       int `gorm:"type:int(10)"`
	CompanyID       int `gorm:"type:int(10) UNSIGNED"`
	Company         Company
	CompanyCategory string `gorm:"type:varchar(255)"`
}

func (sc SublistsCompany) ToEntity() entities.SublistCompany {
	return entities.SublistCompany{
		ID:       strconv.Itoa(int(sc.Company.ID)),
		Name:     sc.Company.CompanyName,
		Category: sc.CompanyCategory,
	}
}

func (s Sublist) TableName() string {
	return "Sublists"
}

func (s Sublist) ToSublistEntity() entities.CompleteSublist {
	return entities.CompleteSublist{
		ID: int(s.ID),
		Location: entities.Location{
			Zipcode: s.Zipcode,
			City:    s.City.Name,
			State:   s.City.State.Name,
			Country: "USA",
		},
		GeneralContractor: entities.GeneralContractorSummary{
			Name: s.WebUser.FullName,
			Company: entities.Company{
				ID:   "1",
				Name: "General Contractor Company",
			},
		},
		ProjectName:     s.ProjectName,
		PublishDate:     s.CreatedAt,
		PublishDateUnix: s.CreatedAt.Unix(),
		Companies:       MakeSublistCompanies(s.SublistsCompanies),
	}
}

func MakeSublistCompanies(sc []SublistsCompany) []entities.SublistCompany {
	esc := make([]entities.SublistCompany, 0)
	for _, s := range sc {
		esc = append(esc, s.ToEntity())
	}
	return esc
}
