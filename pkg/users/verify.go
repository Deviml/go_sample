package users

import (
	"context"
	"fmt"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type VerifyRepository interface {
	VerifyCode(ctx context.Context, userID string, code string, codeType int) error
	VerifyUserSMS(ctx context.Context, userID string) error
	VerifyUserEmail(ctx context.Context, userID string) error
}

type VerifyService struct {
	logger log.Logger
	vr     VerifyRepository
	db     *gorm.DB
	s      *email.Sender
}

func NewVerifyService(logger log.Logger, vr VerifyRepository, db *gorm.DB, s *email.Sender) *VerifyService {
	return &VerifyService{logger: logger, vr: vr, db: db, s: s}
}

func (v VerifyService) Verify(ctx context.Context, userID string, code string, codeType int) error {
	err := v.vr.VerifyCode(ctx, userID, code, codeType)
	if err != nil {
		return err
	}
	if codeType == models.SMSVerificationType {
		err = v.vr.VerifyUserSMS(ctx, userID)
	} else {
		err = v.vr.VerifyUserEmail(ctx, userID)
	}
	var user models.WebUser
	result := v.db.Find(&user, "id = ?", userID)

	if result != nil && result.Error != nil {
		v.logger.Log("find user notify welcome", fmt.Sprintf("%v", err))
		return result.Error
	}

	if user.ProfileType == "vendor_rentals" {

		v.s.SendWelcomeWithCoupon(ctx, user.Username, user.FullName)

	} else {
		err = v.s.SendWelcome(ctx, user.Username, user.FullName)
	}
	if err != nil {
		v.logger.Log("notify welcome", fmt.Sprintf("%v", err))
	}
	return err
}
