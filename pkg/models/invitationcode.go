package models

import "gorm.io/gorm"

type InvitationCode struct {
	gorm.Model
	Used bool
	Code string
}
