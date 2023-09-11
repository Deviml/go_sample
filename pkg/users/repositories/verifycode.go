package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type VerificationCodeRepository struct {
	logger log.Logger
	db     *gorm.DB
}

type CouponRepository struct {
	logger log.Logger
	db     *gorm.DB
}

func (v VerificationCodeRepository) VerifyCode(ctx context.Context, userID string, code string, codeType int) error {
	var verificationCode models.VerificationCode
	result := v.db.Model(&verificationCode).
		Where("web_user_id = ?", userID).
		Where("code = ?", code).
		Where("type = ?", codeType).
		Update("validated", time.Now())
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("not valid code")
	}
	return nil
}

func (v VerificationCodeRepository) VerifyUserSMS(ctx context.Context, userID string) error {
	var user models.WebUser
	result := v.db.Model(&user).Where("id = ?", userID).Update("is_sms_verified", time.Now())
	return result.Error
}

func (v VerificationCodeRepository) VerifyUserEmail(ctx context.Context, userID string) error {
	var user models.WebUser
	result := v.db.Model(&user).Where("id = ?", userID).Update("is_email_verified", time.Now())
	return result.Error
}

func (v VerificationCodeRepository) Create(ctx context.Context, code models.VerificationCode) error {
	result := v.db.Create(&code)
	return result.Error
}

func (c CouponRepository) Create(ctx context.Context, coupon models.Coupon) error {
	result := c.db.Create(&coupon)
	return result.Error
}

func NewVerificationCodeRepository(logger log.Logger, db *gorm.DB) *VerificationCodeRepository {
	return &VerificationCodeRepository{logger: logger, db: db}
}

func (v VerificationCodeRepository) Verify(ctx context.Context, code string) (uint, error) {
	var verificationCode models.InvitationCode
	result := v.db.Where("used = ?", 0).Where("code = ?", code).Find(&verificationCode)
	if result.Error != nil {
		return 0, result.Error
	}
	result = v.db.Model(&verificationCode).Update("used", 1).Where("id = ?", verificationCode.ID)
	if result.Error != nil {
		return 0, result.Error
	}
	return verificationCode.ID, nil
}

func NewCouponRepository(logger log.Logger, db *gorm.DB) *CouponRepository {
	return &CouponRepository{logger: logger, db: db}
}
