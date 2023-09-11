package sublists

import (
	"context"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type UserSublist struct {
	logger log.Logger
	db     *gorm.DB
}

func (us UserSublist) Get(ctx context.Context, userID string) ([]entities.CompleteSublist, error) {
	var s []models.Sublist
	result := us.db.Preload("Location").Preload("City.State").Preload("WebUser").Find(&s, "web_user_id = ?", userID)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	cs := make([]entities.CompleteSublist, 0)
	for _, ms := range s {
		cs = append(cs, ms.ToSublistEntity())
	}
	return cs, nil
}

func NewUserSublist(logger log.Logger, db *gorm.DB) *UserSublist {
	return &UserSublist{logger: logger, db: db}
}
