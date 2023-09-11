package models

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	CouponName string
	Discount   float32
	IsValid    bool
	WebUserID  int
	ExpiredAt  time.Time `gorm:"default:NULL"`
}

func (Coupon) TableName() string {
	return "coupons"
}
