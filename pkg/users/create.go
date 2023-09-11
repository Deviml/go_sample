package users

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/auth"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/go-kit/kit/log"
)

type UserCreateRepository interface {
	CreateVendorRental(ctx context.Context, vendorRental models.VendorRental) (uint, error)
	CreateSubcontractor(ctx context.Context, subcontractor models.Subcontractor) (uint, error)
	CreateGeneralContractor(ctx context.Context, generalContractor models.GeneralContractor, companyName string) (uint, error)
	UpdatePassword(ctx context.Context, email string, password string) error
}

type Hasher interface {
	Hash(element string) (string, error)
}

type VerifyGeneralCode interface {
	Verify(ctx context.Context, code string) (uint, error)
}

type ConfirmationNotifier interface {
	Notify(ctx context.Context, user entities.User, notificationType string) error
	NotifyAgain(ctx context.Context, UserID string, notificationType string) error
	CreateCoupon(ctx context.Context, userID string) error
}

type CreateService struct {
	logger log.Logger
	ucr    UserCreateRepository
	ls     *auth.LoginService
	h      Hasher
	vgc    VerifyGeneralCode
	n      ConfirmationNotifier
	s      email.Sender
}

func (c CreateService) ForgotPassword(ctx context.Context, fp *requests.ForgotPassword) error {
	newPassword := strconv.Itoa(rangeIn(1000000, 9999999))
	hp, err := c.h.Hash(newPassword)
	if err != nil {
		return err
	}
	err = c.ucr.UpdatePassword(ctx, fp.Email, hp)
	if err != nil {
		return err
	}
	return c.s.SendForgotPassword(ctx, fp.Email, newPassword)
}

func NewCreateService(logger log.Logger, ucr UserCreateRepository, ls *auth.LoginService, h Hasher, vgc VerifyGeneralCode, n ConfirmationNotifier, s email.Sender) *CreateService {
	return &CreateService{logger: logger, ucr: ucr, ls: ls, h: h, vgc: vgc, n: n, s: s}
}

func (c CreateService) RequestInvitation(ctx context.Context, ri *requests.RequestInvitation) error {
	return c.s.SendInvitation(ctx, ri.Name, ri.Phone, ri.Email)
}

func (c CreateService) CreateVendorRental(ctx context.Context, user entities.User, profile entities.VendorRental, confirmationType string) (*entities.LoginInformation, error) {
	hashedPassword, err := c.h.Hash(user.Password)
	if err != nil {
		return nil, err
	}
	stateID, _ := strconv.Atoi(profile.StateID)
	cityID, _ := strconv.Atoi(profile.CityID)
	id, err := c.ucr.CreateVendorRental(ctx, models.VendorRental{
		StateID: uint(stateID),
		CityID:  uint(cityID),
		WebUser: models.WebUser{
			Password: hashedPassword,
			Username: user.Email,
			Phone:    user.Phone,
			FullName: user.FullName,
			EveryStateSelection: true,
			Company: models.Company{
				CompanyName: profile.Company.Name,
				Address1:    user.Address,
				Zipcode:     user.Zipcode,
				City:        user.CityName,
				State:       user.StateName,
				Country:     user.CountryName,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	user.ID = strconv.Itoa(int(id))
	return c.createLoginForAccountVendor(ctx, user, confirmationType)
}

func (c CreateService) ResendCode(ctx context.Context, UserID string, confirmationType string) error {
	return c.n.NotifyAgain(ctx, UserID, confirmationType)
}

func (c CreateService) CreateSubcontractor(ctx context.Context, user entities.User, profile entities.Subcontractor, confirmationType string) (*entities.LoginInformation, error) {
	hashedPassword, err := c.h.Hash(user.Password)
	if err != nil {
		return nil, err
	}
	stateID, _ := strconv.Atoi(profile.StateID)
	cityID, _ := strconv.Atoi(profile.CityID)
	id, err := c.ucr.CreateSubcontractor(ctx, models.Subcontractor{
		StateID: uint(stateID),
		CityID:  uint(cityID),
		WebUser: models.WebUser{
			Password: hashedPassword,
			Username: user.Email,
			Phone:    user.Phone,
			FullName: user.FullName,
			Company: models.Company{
				CompanyName: profile.Company.Name,
				Address1:    user.Address,
				Zipcode:     user.Zipcode,
				City:        user.CityName,
				State:       user.StateName,
				Country:     user.CountryName,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	user.ID = strconv.Itoa(int(id))
	return c.createLoginForAccount(ctx, user, confirmationType)
}

func (c CreateService) CreateGeneralContractor(ctx context.Context, user entities.User, profile entities.GeneralContractor, confirmationType string, verificationCode string) (*entities.LoginInformation, error) {
	vID, err := c.vgc.Verify(ctx, verificationCode)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := c.h.Hash(user.Password)
	if err != nil {
		return nil, err
	}
	//err = c.vgc.Verify(ctx, profile.VerificationCode)
	if err != nil {
		return nil, err
	}
	stateID, _ := strconv.Atoi(profile.StateID)
	cityID, _ := strconv.Atoi(profile.CityID)
	id, err := c.ucr.CreateGeneralContractor(ctx, models.GeneralContractor{
		StateID: uint(stateID),
		CityID:  uint(cityID),
		WebUser: models.WebUser{
			Password: hashedPassword,
			Username: user.Email,
			Phone:    user.Phone,
			FullName: user.FullName,
			Company: models.Company{
				CompanyName: profile.Company.Name,
				Address1:    user.Address,
				Zipcode:     user.Zipcode,
				City:        user.CityName,
				State:       user.StateName,
				Country:     user.CountryName,
			},
		},
		InvitationCodeID: &vID,
	}, profile.Company.Name)
	user.ID = strconv.Itoa(int(id))
	if err != nil {
		return nil, err
	}
	return c.createLoginForAccount(ctx, user, confirmationType)
}

func (c CreateService) createLoginForAccount(ctx context.Context, user entities.User, confirmationType string) (*entities.LoginInformation, error) {
	err := c.n.Notify(ctx, user, confirmationType)
	if err != nil {
		c.logger.Log("notify", fmt.Sprintf("%v", err))
		return nil, err
	}

	c.logger.Log("user", fmt.Sprintf("%v", user))
	return c.ls.Login(ctx, user.Email, user.Password)
}

func (c CreateService) createLoginForAccountVendor(ctx context.Context, user entities.User, confirmationType string) (*entities.LoginInformation, error) {
	err := c.n.Notify(ctx, user, confirmationType)
	if err != nil {
		c.logger.Log("notify", fmt.Sprintf("%v", err))
		return nil, err
	}

	create := c.n.CreateCoupon(ctx, user.ID)
	if create != nil {
		c.logger.Log("coupon", fmt.Sprintf("%v", create))
		return nil, create
	}

	c.logger.Log("user", fmt.Sprintf("%v", user))
	return c.ls.Login(ctx, user.Email, user.Password)
}
