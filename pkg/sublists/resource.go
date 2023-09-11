package sublists

import (
	"context"
	"fmt"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/email"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/endpoints/sublists"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
	"strconv"
)

type Services struct {
	logger log.Logger
	db     *gorm.DB
	es     email.Sender
}

func NewServices(logger log.Logger, db *gorm.DB, es email.Sender) *Services {
	return &Services{logger: logger, db: db, es: es}
}

func (s Services) CreateSublistForUser(ctx context.Context, userID string, createRequest sublists.CreateRequest) error {
	//strconv.Atoi(createRequest.LocationID)
	locationID := 1
	uID, _ := strconv.Atoi(userID)
	cID, _ := strconv.Atoi(createRequest.CityID)
	_, _ = strconv.Atoi(createRequest.CountyID)
	sublist := models.Sublist{
		ProjectName: createRequest.Name,
		LocationID:  locationID,
		WebUserID:   uID,
		CityID:      cID,
		CountyID:    nil,
		Zipcode:     createRequest.Zipcode,
	}
	result := s.db.Create(&sublist)
	if result != nil && result.Error != nil {
		return result.Error
	}

	for _, c := range createRequest.Companies {
		if c.ID == "" && c.Name == "" {
			continue
		}
		var company models.Company
		var subcompany models.SublistsCompany
		if c.ID == "" {
			company = models.Company{
				CompanyName: c.Name,
			}
			subcompany = models.SublistsCompany{
				SublistID:       int(sublist.ID),
				Company:         company,
				CompanyCategory: c.Category,
			}
		} else {
			id, _ := strconv.Atoi(c.ID)
			subcompany = models.SublistsCompany{
				SublistID:       int(sublist.ID),
				CompanyID:       id,
				CompanyCategory: c.Category,
			}
		}

		result = s.db.Create(&subcompany)
	}
	result = s.db.Preload("Location").Preload("City.State").Find(&sublist, "id = ?", sublist.ID)
	if result.Error != nil {
		s.logger.Log(fmt.Sprintf("err sublist notify create %s", result.Error.Error()))
		return nil
	}
	err := s.notifyCreate(ctx, sublist)
	if err != nil {
		return s.logger.Log("err", err.Error())
	}
	return nil
}

func (s Services) notifyCreate(ctx context.Context, sm models.Sublist) error {
	var users []models.WebUser
	result := s.db.Preload("Cities").Preload("Counties").Preload("States").Find(&users, "profile_type = ?", "vendor_rentals")
	if result.Error != nil {
		return result.Error
	}
	emails := make([]string, 0)
	for _, u := range users {
		if !u.SublistNotification {
			continue
		}
		s.logger.Log("user_states", fmt.Sprintf("%v", u.States))
		s.logger.Log("state_id", sm.City.StateID)
		s.logger.Log("city_id", sm.City.ID)

		if len(u.States) > 0 {
			if len(u.Cities) > 0 {
				if !u.HaveCity(uint(sm.CityID)) {
					s.logger.Log("skip_user", u.ID)
					continue
				}
			}
			if !u.HaveState(sm.City.StateID) {
				s.logger.Log("skip_user", u.ID)
				continue
			}
		} else {
			s.logger.Log("skip_user_no_preference", u.ID)
			continue
		}

		emails = append(emails, u.Username)
	}
	s.logger.Log("email", fmt.Sprintf("%v", emails))
	if len(emails) == 0 {
		return nil
	}
	return s.es.SendNewSublist(ctx, emails, []models.Sublist{sm})
}

func (s Services) UpdateSublist(ctx context.Context, updateRequest sublists.UpdateRequest) error {
	var sublist models.Sublist
	result := s.db.Preload("City.State").Preload("WebUsers").Find(&sublist, "id = ?", updateRequest.ID)
	if result != nil && result.Error != nil {
		return result.Error
	}

	result = s.db.Model(&models.Sublist{}).Where("id = ?", updateRequest.ID).Update("zipcode", updateRequest.Zipcode)
	if result != nil && result.Error != nil {
		return result.Error
	}

	if updateRequest.Name != "" {
		result = s.db.Model(&models.Sublist{}).Where("id = ?", updateRequest.ID).Update("Project_Name", updateRequest.Name)
		if result != nil && result.Error != nil {
			return result.Error
		}
	}

	var sc models.SublistsCompany
	result = s.db.Delete(&sc, "sublist_id = ?", updateRequest.ID)
	scs := make([]models.SublistsCompany, 0)
	for _, c := range updateRequest.Companies {
		if c.ID == "" && c.Name == "" {
			continue
		}
		var company models.Company
		var subcompany models.SublistsCompany
		if c.ID == "" {
			company = models.Company{
				CompanyName: c.Name,
			}
			subcompany = models.SublistsCompany{
				SublistID:       int(sublist.ID),
				Company:         company,
				CompanyCategory: c.Category,
			}
		} else {
			id, _ := strconv.Atoi(c.ID)
			subcompany = models.SublistsCompany{
				SublistID:       int(sublist.ID),
				CompanyID:       id,
				CompanyCategory: c.Category,
			}
		}

		result = s.db.Create(&subcompany)
		if result.Error != nil {
			s.logger.Log("err", result.Error)
		} else {
			scs = append(scs, subcompany)
		}
	}
	err := s.notifyChange(ctx, sublist, scs)
	if err != nil {
		s.logger.Log(fmt.Sprintf("email error: %s", err.Error()))
	}
	return nil
}

func (s Services) notifyChange(ctx context.Context, sm models.Sublist, sc []models.SublistsCompany) error {
	destinations := make([]string, 0)
	for _, u := range sm.WebUsers {
		destinations = append(destinations, u.Username)
	}
	return s.es.SendSublistChange(ctx, destinations, sm, sc)
}

func (s Services) ShowSublist(ctx context.Context, sublistID string) (entities.CompleteSublist, error) {
	var sublist models.Sublist
	result := s.db.Preload("Location").Preload("City.State").Preload("WebUser").Find(&sublist, "id = ?", sublistID)
	if result.Error != nil {
		return entities.CompleteSublist{}, result.Error
	}
	var sc []models.SublistsCompany
	result = s.db.Preload("Company").Find(&sc, "sublist_id = ?", sublistID)
	if result.Error != nil {
		return entities.CompleteSublist{}, result.Error
	}
	sublist.SublistsCompanies = sc
	return sublist.ToSublistEntity(), nil
}

func (s Services) DeleteSublist(ctx context.Context, sublistID string) error {
	var sublist models.Sublist
	result := s.db.Delete(&sublist, "id = ?", sublistID)
	if result != nil && result.Error != nil {
		return result.Error
	}
	return nil
}
