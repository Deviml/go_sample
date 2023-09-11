package auth

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"strconv"
)

type CredentialsIssuer interface {
	Issue(user entities.User) (*entities.LoginInformation, error)
}

type HashComparator interface {
	Compare(expected string, actual string) error
}

type UserRepository interface {
	FindUserByUserName(ctx context.Context, username string) (*models.WebUser, error)
}

type LoginService struct {
	userRepository UserRepository
	comparator     HashComparator
	issuer         CredentialsIssuer
}

func NewLoginService(userRepository UserRepository, hasher HashComparator, authorizer CredentialsIssuer) *LoginService {
	return &LoginService{userRepository: userRepository, comparator: hasher, issuer: authorizer}
}

func (l LoginService) Login(ctx context.Context, username string, password string) (*entities.LoginInformation, error) {
	user, err := l.userRepository.FindUserByUserName(ctx, username)
	if err != nil {
		return nil, err
	}
	err = l.comparator.Compare(user.Password, password)
	if err != nil {
		return nil, err
	}
	return l.issuer.Issue(entities.User{
		ID:       strconv.Itoa(int(user.ID)),
		Password: user.Password,
	})
}
