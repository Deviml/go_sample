package users

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

const (
	SMSNotificationType   = "1"
	EmailNotificationType = "2"
)

type CreateVerificationCodeRepository interface {
	Create(ctx context.Context, code models.VerificationCode) error
}

type CreateCouponRepository interface {
	Create(ctx context.Context, coupon models.Coupon) error
}

func rangeIn(low, hi int) int {
	rand.Seed(time.Now().UnixNano())
	return low + rand.Intn(hi-low)
}

type UserNotifier struct {
	logger log.Logger
	db     *gorm.DB
	cvcr   CreateVerificationCodeRepository
	ccr    CreateCouponRepository
	snsC   *sns.Client
	s      email.Sender
}

func NewUserNotifier(logger log.Logger, db *gorm.DB, cvcr CreateVerificationCodeRepository, ccr CreateCouponRepository, snsC *sns.Client, s email.Sender) *UserNotifier {
	return &UserNotifier{logger: logger, db: db, cvcr: cvcr, ccr: ccr, snsC: snsC, s: s}
}

func (u UserNotifier) CreateCoupon(ctx context.Context, userID string) error {
	uID, _ := strconv.Atoi(userID)
	coupon := models.Coupon{
		CouponName: "EQ50",
		Discount:   50,
		IsValid:    true,
		WebUserID:  uID,
	}
	err := u.ccr.Create(ctx, coupon)
	return err
}

func (u UserNotifier) NotifyAgain(ctx context.Context, UserID string, notificationType string) error {

	var user models.WebUser
	uID, _ := strconv.Atoi(UserID)
	err := u.db.Where("id = ?", UserID).First(&user).Error
	if err != nil {
		return err
	}

	if notificationType == EmailNotificationType {
		verificationCode := models.VerificationCode{
			Code:      strconv.Itoa(rangeIn(1000, 9999)),
			Type:      models.EmailVerificationType,
			WebUserID: uint(uID),
		}
		err := u.cvcr.Create(ctx, verificationCode)
		if err != nil {
			return err
		}
		return u.SendConfirmationEmailAgain(ctx, user.Username, verificationCode)
	}
	return nil
}

func (u UserNotifier) Notify(ctx context.Context, user entities.User, notificationType string) error {
	uID, _ := strconv.Atoi(user.ID)
	if notificationType == EmailNotificationType {
		verificationCode := models.VerificationCode{
			Code:      strconv.Itoa(rangeIn(1000, 9999)),
			Type:      models.EmailVerificationType,
			WebUserID: uint(uID),
		}
		err := u.cvcr.Create(ctx, verificationCode)
		if err != nil {
			return err
		}
		return u.SendConfirmationEmail(ctx, user, verificationCode)
	}

	verificationCode := models.VerificationCode{
		Code:      strconv.Itoa(rangeIn(1000, 9999)),
		Type:      models.SMSVerificationType,
		WebUserID: uint(uID),
	}
	err := u.cvcr.Create(ctx, verificationCode)
	if err != nil {
		u.logger.Log("create", fmt.Sprintf("create:%v", err))
		return nil
	}
	err = u.SendConfirmationSMS(ctx, user, verificationCode)
	return err
}

func (u UserNotifier) SendConfirmationEmail(ctx context.Context, user entities.User, code models.VerificationCode) error {
	return u.s.SendConfirmation(ctx, user.Email, code.Code)
}

func (u UserNotifier) SendConfirmationEmailAgain(ctx context.Context, Email string, code models.VerificationCode) error {
	return u.s.SendConfirmation(ctx, Email, code.Code)
}

func (u UserNotifier) SendConfirmationSMS(ctx context.Context, user entities.User, code models.VerificationCode) error {
	message := fmt.Sprintf("Hello, %s \n This is your confirmation code for Equiphunter account: %s", user.FullName, code.Code)
	input := &sns.PublishInput{
		Message:     aws.String(message),
		PhoneNumber: aws.String(fmt.Sprintf("+1%s", user.Phone)),
	}
	req := u.snsC.PublishRequest(input)
	resp, err := req.Send(ctx)
	if err != nil {
		u.logger.Log("sns", fmt.Sprintf("ses:%v", err))
		return err
	}
	u.logger.Log("sns", fmt.Sprintf("ses:%v", resp))
	return nil
}
