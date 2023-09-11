package users

import (
	"context"
	"strconv"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type ResourceService struct {
	logger log.Logger
	db     *gorm.DB
	h      Hasher
}

func NewResourceService(logger log.Logger, db *gorm.DB, h Hasher) *ResourceService {
	return &ResourceService{logger: logger, db: db, h: h}
}

func (r ResourceService) PatchUser(ctx context.Context, userID string, request *requests.UpdateUserRequest) error {
	uid, _ := strconv.Atoi(userID)
	user := models.WebUser{
		Model: gorm.Model{
			ID: uint(uid),
		},
	}
	result := r.db.Model(&user).Where("id = ?", userID).
		Update("sublist_notification", request.SublistNotification).
		Update("equipment_request_notification", request.EquipmentRequestNotification).
		Update("supply_request_notification", request.SupplyRequestNotification).
		Update("new_user_pop_up", request.NewUserPopUp).
		Update("every_state_selection", request.EveryStateSelection)

	if result.Error != nil {
		return result.Error
	}

	if request.Password != "" {
		hPassword, _ := r.h.Hash(request.Password)
		result := r.db.Model(&user).Where("id = ?", userID).
			Update("password", hPassword)
		if result.Error != nil {
			return result.Error
		}
	}

	supplies := make([]models.SupplyCategory, 0)
	for _, i := range request.SupplyCategories {
		supplies = append(supplies, models.SupplyCategory{
			ID: i,
		})
	}

	err := r.db.Model(&user).Association("SupplyCategories").Replace(supplies)
	if err != nil {
		return err
	}

	equipments := make([]models.EquipmentCategory, 0)
	for _, i := range request.EquipmentCategories {
		equipments = append(equipments, models.EquipmentCategory{
			ID: i,
		})
	}

	err = r.db.Model(&user).Association("EquipmentCategories").Replace(equipments)
	if err != nil {
		return err
	}

	cities := make([]models.City, 0)
	for _, i := range request.Cities {
		cities = append(cities, models.City{
			Model: gorm.Model{ID: i},
		})
	}
	err = r.db.Model(&user).Association("Cities").Replace(cities)
	if err != nil {
		return err
	}

	states := make([]models.State, 0)
	for _, i := range request.States {
		states = append(states, models.State{
			Model: gorm.Model{
				ID: i,
			},
		})
	}
	err = r.db.Model(&user).Association("States").Replace(states)

	if err != nil {
		return err
	}

	counties := make([]models.County, 0)
	for _, i := range request.Counties {
		counties = append(counties, models.County{
			Model: gorm.Model{
				ID: i,
			},
		})
	}
	err = r.db.Model(&user).Association("Counties").Replace(counties)
	return nil
}

func (r ResourceService) GetUser(ctx context.Context, userID string) (*entities.UserInfo, error) {
	var user models.WebUser
	result := r.db.Preload("VerificationCodes").Find(&user, "id = ?", userID)
	if result.Error != nil {
		return nil, result.Error
	}

	var supplies []models.SupplyCategory
	err := r.db.Model(&user).Association("SupplyCategories").Find(&supplies)
	if err != nil {
		return nil, err
	}
	user.SupplyCategories = supplies
	var equipments []models.EquipmentCategory
	err = r.db.Model(&user).Association("EquipmentCategories").Find(&equipments)
	if err != nil {
		return nil, err
	}
	user.EquipmentCategories = equipments

	var cities []models.City
	err = r.db.Model(&user).Association("Cities").Find(&cities)
	if err != nil {
		return nil, err
	}
	user.Cities = cities

	var counties []models.County
	err = r.db.Model(&user).Association("Counties").Find(&counties)
	if err != nil {
		return nil, err
	}
	user.Counties = counties

	var states []models.State
	err = r.db.Model(&user).Association("States").Find(&states)
	if err != nil {
		return nil, err
	}
	user.States = states
	userInfo := user.ToUserInfo()
	return &userInfo, nil
}

func (r ResourceService) CloseAccount(ctx context.Context, userID string) error {

	var user models.WebUser
	result := r.db.Find(&user, "id = ?", userID)
	if result.Error != nil {
		return result.Error
	}
	result = r.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r ResourceService) UpdateAccount(ctx context.Context, userID string, request *requests.UpdateAccount) error {
	var user models.WebUser

	if request.FullName != "" {
		result := r.db.Model(&user).Where("id = ?", userID).
			Update("full_name", request.FullName)
		if result.Error != nil {
			return result.Error
		}
	}
	if request.Phone != "" {
		result := r.db.Model(&user).Where("id = ?", userID).
			Update("phone", request.Phone)
		if result.Error != nil {
			return result.Error
		}
	}
	if request.ProfilePicture != "" {
		result := r.db.Model(&user).Where("id = ?", userID).
			Update("profile_picture", request.ProfilePicture)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
