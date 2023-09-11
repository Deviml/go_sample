package models

import (
	"strconv"
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"gorm.io/gorm"
)

type Company struct {
	ID          uint `gorm:"primarykey, type:int(10) UNSIGNED"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	CompanyName string         `gorm:"type:varchar(255)"`
	Address1    string         `gorm:"type:varchar(255)"`
	Zipcode     string         `gorm:"type:varchar(255)"`
	City        string         `gorm:"type:varchar(255)"`
	State       string         `gorm:"type:varchar(255)"`
	Country     string         `gorm:"type:varchar(255)"`
	WebUserID   uint           `gorm:"type:int(10) UNSIGNED"`
}

func (Company) TableName() string {
	return "companies"
}

func (c Company) ToEntity() entities.Company {
	return entities.Company{
		ID:   strconv.Itoa(int(c.ID)),
		Name: c.CompanyName,
	}
}
