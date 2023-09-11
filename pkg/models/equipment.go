package models

import (
	"strconv"
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"gorm.io/gorm"
)

type EquipmentCategory struct {
	ID        uint `gorm:"primarykey;type:int(11) UNSIGNED"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
}

func (e EquipmentCategory) TableName() string {
	return "EquipmentCategories"
}

func (e EquipmentCategory) ToEntity() entities.EquipmentCategory {
	return entities.EquipmentCategory{
		ID:   strconv.Itoa(int(e.ID)),
		Name: e.Name,
	}
}

type EquipmentSubcategory struct {
	gorm.Model
	Name                string `gorm:"type:varchar(255)"`
	EquipmentCategoryID uint   `gorm:"type:int(11)"`
	EquipmentCategory   EquipmentCategory
}

func (e EquipmentSubcategory) TableName() string {
	return "EquipmentSubcategories"
}

func (e EquipmentSubcategory) ToEntity() entities.EquipmentSubcategory {
	return entities.EquipmentSubcategory{
		ID:                strconv.Itoa(int(e.ID)),
		Name:              e.Name,
		EquipmentCategory: e.EquipmentCategory.ToEntity(),
	}
}

type Equipment struct {
	gorm.Model
	Name                   string
	EquipmentSubcategoryID uint `gorm:"type:int(11)"`
	EquipmentSubcategory   EquipmentSubcategory
}

func (e Equipment) TableName() string {
	return "equipments"
}

func (e Equipment) ToEntity() entities.Equipment {
	return entities.Equipment{
		ID:          strconv.Itoa(int(e.ID)),
		Name:        e.Name,
		Subcategory: e.EquipmentSubcategory.ToEntity(),
	}
}
