package mysql

import (
	"context"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type User struct {
	logger log.Logger
	db     *gorm.DB
}

func NewUser(logger log.Logger, db *gorm.DB) *User {
	return &User{logger: logger, db: db}
}

func (u User) FindUserByUserName(ctx context.Context, username string) (*models.WebUser, error) {
	configuration := &models.WebUser{}
	if result := u.db.First(configuration, "username = ?", username); result.Error != nil {
		return nil, result.Error
	}
	return configuration, nil
}
