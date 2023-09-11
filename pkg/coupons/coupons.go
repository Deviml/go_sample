package coupons

import (
	"strconv"
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	coupons2 "github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/coupons"
	"github.com/go-kit/kit/log"

	"gorm.io/gorm"
)

type Service struct {
	logger log.Logger
	db     *gorm.DB
}

func NewService(logger log.Logger, db *gorm.DB) *Service {
	return &Service{logger: logger, db: db}
}

func (s *Service) ValidateCoupon(req coupons2.ValidateCouponRequest, userID string) (entities.Coupon, error) {

	uid, _ := strconv.Atoi(userID)

	var coupon models.Coupon
	err := s.db.Where("coupon_name = ? AND web_user_id = ? AND is_valid = ?", req.CouponName, uid, true).First(&coupon).Error

	if err != nil {
		return entities.Coupon{Message: "Coupon does not exist or has expired!"}, nil
	}

	return entities.Coupon{CouponName: coupon.CouponName, Discount: coupon.Discount, IsValid: coupon.IsValid, UserID: userID, Message: "Coupon is valid!"}, nil
}

func (s *Service) UpdateCoupon(req coupons2.UpdateCouponRequest, userID string) (string, error) {

	uid, _ := strconv.Atoi(userID)

	var coupon models.Coupon
	err := s.db.Where("coupon_name = ? AND web_user_id = ? AND is_valid = ?", req.CouponName, uid, true).First(&coupon).Error

	if err != nil {
		return "Coupon does not exist!", nil
	}

	coupon.IsValid = false
	coupon.ExpiredAt = time.Now()
	s.db.Save(&coupon)

	return "Coupon updated!", nil
}
