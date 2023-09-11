package claims

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

func NewUserClaims(userID string, expirationTime int64) *UserClaims {
	return &UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  time.Now().Unix(),
		},
		UserID: userID,
	}
}
