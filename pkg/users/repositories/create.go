package repositories

import (
	"context"
	"errors"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type Create struct {
	logger log.Logger
	db     *gorm.DB
}

func (c Create) UpdatePassword(ctx context.Context, email string, password string) error {
	var fUser models.WebUser
	result := c.db.Find(&fUser, "username = ?", email)
	if result != nil && result.Error != nil {
		return result.Error
	}
	if fUser.ID == 0 {
		return errors.New("user not found")
	}
	var user models.WebUser
	result = c.db.Model(&user).Where("username = ?", email).Update("password", password)
	return result.Error
}

func NewCreate(logger log.Logger, db *gorm.DB) *Create {
	return &Create{logger: logger, db: db}
}

func (c Create) CreateVendorRental(ctx context.Context, vendorRental models.VendorRental) (uint, error) {
	if err := c.ValidateUserExists(ctx, vendorRental.WebUser.Username); err != nil {
		return 0, err
	}
	result := c.db.Create(&vendorRental)
	return vendorRental.WebUser.ID, result.Error
}

func (c Create) CreateSubcontractor(ctx context.Context, subcontractor models.Subcontractor) (uint, error) {
	if err := c.ValidateUserExists(ctx, subcontractor.WebUser.Username); err != nil {
		return 0, err
	}
	result := c.db.Create(&subcontractor)
	return subcontractor.WebUser.ID, result.Error
}

func (c Create) CreateGeneralContractor(ctx context.Context, generalContractor models.GeneralContractor, companyName string) (uint, error) {
	if err := c.ValidateUserExists(ctx, generalContractor.WebUser.Username); err != nil {
		return 0, err
	}
	if generalContractor.CompanyID == 0 {
		var cp models.Company
		cp.CompanyName = companyName
		c.db.Create(&cp)
		generalContractor.CompanyID = cp.ID
	}
	result := c.db.Create(&generalContractor)
	return generalContractor.WebUser.ID, result.Error
}

func (c Create) ValidateUserExists(ctx context.Context, username string) error {
	var wu models.WebUser
	result := c.db.Find(&wu, "username = ?", username)
	if result != nil && result.Error == nil && wu.ID != 0 {
		return errors.New("user exists")
	}
	return nil
}
