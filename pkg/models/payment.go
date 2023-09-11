package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	PayedStatus  = 1
	RefundStatus = 2
)

type Payment struct {
	gorm.Model
	Payload        datatypes.JSON
	Response       datatypes.JSON
	WebUserID      int
	WebUser        WebUser
	Total          float32
	Status         uint
	Email          string
	NameOnCard     string
	Zipcode        string
	PaymentDetails []PaymentDetail
}

type PaymentDetail struct {
	gorm.Model
	PaymentID int
	Payment   Payment
	QuoteID   *int `gorm:"type:int(10)"`
	SublistID *int `gorm:"type:int(10)"`
}
