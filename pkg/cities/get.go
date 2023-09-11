package cities

import (
	"context"
	"fmt"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/helpers"
	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
)

type GetCitiesService struct {
	logger log.Logger
	db     *gorm.DB
}

func (g GetCitiesService) ListState(ctx context.Context, keywords []string) ([]entities.State, error) {
	var modelStates []models.State
	query := g.db
	for idx, keyword := range keywords {
		if idx == 0 {
			query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", keyword))
			continue
		}
		query = query.Or("name LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	result := query.Find(&modelStates)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var states []entities.State
	for _, modelState := range modelStates {
		states = append(states, modelState.ToEntity())
	}
	states = helpers.SortStatesWithoutRepeating(states)
	return states, nil
}

func (g GetCitiesService) ListCounties(ctx context.Context, keywords []string, stateID string) ([]entities.County, error) {
	var modelCounties []models.County
	query := g.db
	for idx, keyword := range keywords {
		if idx == 0 {
			query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", keyword))
			continue
		}
		query = query.Or("name LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	//if len(keywords) == 0 {
	//	query = query.Limit(5)
	//}
	result := query.Distinct("name", "id")
	if stateID != "" {
		result = result.Where("state_id = ?", stateID)
	}
	result = result.Find(&modelCounties)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var counties []entities.County
	for _, modelCounty := range modelCounties {
		counties = append(counties, modelCounty.ToEntity())
	}
	counties = helpers.SortCountiesWithoutRepeating(counties)
	counties = helpers.ConvertToLowercase(counties)
	counties = helpers.CapitalizeFirstLetterOfCounty(counties)
	return counties, nil
}

func NewGetCitiesService(logger log.Logger, db *gorm.DB) *GetCitiesService {
	return &GetCitiesService{logger: logger, db: db}
}

func (g GetCitiesService) ListCities(ctx context.Context, keywords []string, stateID string) ([]entities.City, error) {
	var modelCities []models.City
	query := g.db
	for idx, keyword := range keywords {
		if idx == 0 {
			query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", keyword))
			continue
		}
		query = query.Or("name LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}

	if stateID != "" {
		query = query.Where("state_id = ?", stateID)
	}

	result := query.Find(&modelCities)
	if result != nil && result.Error != nil {
		return nil, result.Error
	}
	var cities []entities.City
	for _, modelCity := range modelCities {
		cities = append(cities, modelCity.ToEntity())
	}
	cities = helpers.SortCitiesWithoutRepeating(cities)
	return cities, nil
}
