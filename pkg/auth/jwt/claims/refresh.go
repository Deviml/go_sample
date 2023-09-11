package claims

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Refresh struct {
	jwt.StandardClaims
	Token string `json:"token"`
}

func NewRefresh(token string, expiration int64) Refresh {
	return Refresh{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration,
			IssuedAt:  time.Now().Unix(),
		},
		Token: token,
	}
}
