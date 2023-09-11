package hmac

import (
	"context"
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/auth/jwt/claims"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/dgrijalva/jwt-go"
)

type Issuer struct {
	secretKey         []byte
	expirationTime    int64
	refreshExpiration int64
}

func (i Issuer) IssueRefresh(ctx context.Context, userID string) (*entities.LoginInformation, error) {
	return i.makeCredentials(userID)
}

func NewIssuer(secretKey []byte, expirationTime int64, refreshExpiration int64) *Issuer {
	return &Issuer{secretKey: secretKey, expirationTime: expirationTime, refreshExpiration: refreshExpiration}
}

func (i Issuer) Issue(user entities.User) (*entities.LoginInformation, error) {
	return i.makeCredentials(user.ID)
}

func (i Issuer) makeCredentials(userID string) (*entities.LoginInformation, error) {
	userClaims := claims.NewUserClaims(userID, time.Now().Unix()+i.expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)

	tokenString, err := token.SignedString(i.secretKey)
	if err != nil {
		return nil, err
	}
	refreshClaims := claims.NewRefresh(tokenString, time.Now().Unix()+i.refreshExpiration)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refresh.SignedString(i.secretKey)
	if err != nil {
		return nil, err
	}
	return &entities.LoginInformation{
		Token:        tokenString,
		RefreshToken: refreshString,
		Claims:       *userClaims,
		UserID:       userID,
	}, nil
}
