package models

import (
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type SupplyCategory struct {
	ID        uint `gorm:"primarykey;type:int(11) UNSIGNED"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
}

func (sc SupplyCategory) ToEntity() entities.SupplyCategory {
	return entities.SupplyCategory{
		ID:   strconv.Itoa(int(sc.ID)),
		Name: sc.Name,
	}
}

func (sc SupplyCategory) TableName() string {
	return "SupplyCategories"
}

type Supply struct {
	gorm.Model
	Name             string
	SupplyCategoryID uint
	SupplyCategory   SupplyCategory
}

func (s Supply) ToEntity() entities.Supply {
	return entities.Supply{
		ID:             strconv.Itoa(int(s.ID)),
		Name:           s.Name,
		SupplyCategory: s.SupplyCategory.ToEntity(),
	}
}
