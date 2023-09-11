package models

import (
	"database/sql"
	"strconv"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"gorm.io/gorm"
)

const (
	EmailVerificationType = 1
	SMSVerificationType   = 2
	GeneralContractorType = "general_contractors"
	VendorRentalsType     = "vendor_rentals"
	SubcontractorType     = "subcontractors"
	GCEnum                = 1
	VREnum                = 3
	SCEnum                = 2
	VerificationTypeEmail = "1"
	VerificationTypeSMS   = "2"
)

type VerificationCode struct {
	gorm.Model
	Code      string
	Validated sql.NullTime
	Type      int
	WebUserID uint
	WebUser   WebUser
}

type WebUser struct {
	gorm.Model
	Password                     string `gorm:"varchar(50)"`
	Username                     string `gorm:"varchar(250)"`
	Phone                        string `gorm:"varchar(250)"`
	FullName                     string `gorm:"varchar(250)"`
	ProfilePicture               string
	ProfileID                    string
	ProfileType                  string
	IsSMSVerified                sql.NullTime
	IsEmailVerified              sql.NullTime
	VerificationCodes            []VerificationCode
	SublistNotification          bool
	EquipmentRequestNotification bool
	SupplyRequestNotification    bool
	EveryStateSelection          bool `gorm:"default:true"`
	NewUserPopUp                 bool
	Company                      Company
	Coupon                       Coupon
	SupplyCategories             []SupplyCategory    `gorm:"many2many:user_supply_categories;"`
	EquipmentCategories          []EquipmentCategory `gorm:"many2many:user_equipment_categories;"`
	Sublists                     []Sublist
	Cities                       []City   `gorm:"many2many:user_cities;"`
	Counties                     []County `gorm:"many2many:user_counties;"`
	States                       []State  `gorm:"many2many:user_states;"`
	Payments                     []Payment
}

func (w WebUser) ToUserInfo() entities.UserInfo {
	supplies := make([]uint, 0)
	for _, s := range w.SupplyCategories {
		supplies = append(supplies, s.ID)
	}
	equipments := make([]uint, 0)
	for _, e := range w.EquipmentCategories {
		equipments = append(equipments, e.ID)
	}

	cities := make([]entities.City, 0)
	for _, c := range w.Cities {
		cities = append(cities, entities.City{
			ID:   strconv.Itoa(int(c.ID)),
			Name: c.Name,
		})
	}

	states := make([]entities.State, 0)
	for _, s := range w.States {
		states = append(states, entities.State{
			ID:   strconv.Itoa(int(s.ID)),
			Name: s.Name,
		})
	}

	counties := make([]entities.County, 0)
	for _, c := range w.Counties {
		counties = append(counties, entities.County{
			ID:   strconv.Itoa(int(c.ID)),
			Name: c.Name,
		})
	}

	verificationType := VerificationTypeEmail

	if len(w.VerificationCodes) > 0 {
		verificationType = strconv.Itoa(w.VerificationCodes[0].Type)
	}

	return entities.UserInfo{
		ID:                           int(w.ID),
		AccountType:                  w.GetUserType(),
		SublistNotification:          w.SublistNotification,
		EquipmentRequestNotification: w.EquipmentRequestNotification,
		SupplyRequestNotification:    w.SupplyRequestNotification,
		EveryStateSelection:          w.EveryStateSelection,
		NewUserPopUp:                 w.NewUserPopUp,
		SupplyCategories:             supplies,
		EquipmentCategories:          equipments,
		Cities:                       cities,
		Counties:                     counties,
		States:                       states,
		IsVerified:                   w.IsEmailVerified.Valid || w.IsSMSVerified.Valid,
		VerificationType:             verificationType,
		FullName:                     w.FullName,
		PhoneNumber:                  w.Phone,
		Email:                        w.Username,
		ProfilePicture:               w.ProfilePicture,
	}
}

func (w WebUser) GetUserType() int {
	switch w.ProfileType {
	case GeneralContractorType:
		return GCEnum
	case SubcontractorType:
		return SCEnum
	case VendorRentalsType:
		return VREnum
	}
	return VREnum
}

func (w WebUser) HaveCity(cID uint) bool {
	for _, c := range w.Cities {
		if c.ID == cID {
			return true
		}
	}
	return false
}

func (w WebUser) HaveState(sID uint) bool {
	for _, s := range w.States {
		if s.ID == sID {
			return true
		}
	}
	return false
}

func (w WebUser) HaveCounty(cID uint) bool {
	for _, c := range w.Counties {
		if c.ID == cID {
			return true
		}
	}
	return false
}

func (w WebUser) HaveEquipmentCategory(eqID uint) bool {
	for _, eq := range w.EquipmentCategories {
		if eq.ID == eqID {
			return true
		}
	}
	return false
}

func (w WebUser) HaveSupplyCategory(sqID uint) bool {
	for _, sq := range w.SupplyCategories {
		if sq.ID == sqID {
			return true
		}
	}
	return false
}

type GeneralContractor struct {
	gorm.Model
	CompanyID        uint
	StateID          uint
	CityID           uint
	WebUser          WebUser `gorm:"polymorphic:Profile;"`
	InvitationCodeID *uint
	InvitationCode   InvitationCode
}

type Subcontractor struct {
	gorm.Model
	StateID uint
	CityID  uint
	WebUser WebUser `gorm:"polymorphic:Profile;"`
}

type VendorRental struct {
	gorm.Model
	StateID uint
	CityID  uint
	WebUser WebUser `gorm:"polymorphic:Profile;"`
	Quotes  []Quote `gorm:"many2many:vendor_rental_quotes;"`
}
