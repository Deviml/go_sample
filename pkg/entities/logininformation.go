package entities

import "github.com/Equiphunter-com/equipment-hunter-api/pkg/auth/jwt/claims"

type LoginInformation struct {
	Token        string
	RefreshToken string
	Claims       claims.UserClaims
	UserID       string
}
