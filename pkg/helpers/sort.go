package helpers

import (
	"sort"
	"strings"
	"time"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/entities"
)

func SortRequestsByPurchaseDateDesc(requests []entities.CompleteEquipmentRequest) []entities.CompleteEquipmentRequest {
	sort.SliceStable(requests, func(i, j int) bool {
		return requests[i].PurchasedAt > requests[j].PurchasedAt
	})
	return requests
}

func SortRequestsByPurchaseDateAsc(requests []entities.CompleteEquipmentRequest) []entities.CompleteEquipmentRequest {
	sort.SliceStable(requests, func(i, j int) bool {
		return requests[i].PurchasedAt < requests[j].PurchasedAt
	})
	return requests
}

func SortRequestsByRequestDateAsc(requests []entities.CompleteEquipmentRequest) []entities.CompleteEquipmentRequest {
	sort.SliceStable(requests, func(i, j int) bool {
		return requests[i].StartDate < requests[j].StartDate
	})
	return requests
}

func SortRequestsByRequestDateDesc(requests []entities.CompleteEquipmentRequest) []entities.CompleteEquipmentRequest {
	sort.SliceStable(requests, func(i, j int) bool {
		return requests[i].StartDate > requests[j].StartDate
	})
	return requests
}

func SortQuotesByStartDate(quotes []entities.CompleteQuote) []entities.CompleteQuote {
	sort.SliceStable(quotes, func(i, j int) bool {
		return quotes[i].StartDate > quotes[j].StartDate
	})
	return quotes
}

func SortQuotesByLocationAsc(quotes []entities.Quote) []entities.Quote {
	sort.SliceStable(quotes, func(i, j int) bool {
		return strings.ToLower(quotes[i].Location.City) < strings.ToLower(quotes[j].Location.City)
	})
	return quotes
}

func SortQuotesByLocationDesc(quotes []entities.Quote) []entities.Quote {
	sort.SliceStable(quotes, func(i, j int) bool {
		return strings.ToLower(quotes[i].Location.City) > strings.ToLower(quotes[j].Location.City)
	})
	return quotes
}

func SortQuotesByNameAsc(quotes []entities.Quote) []entities.Quote {
	sort.SliceStable(quotes, func(i, j int) bool {
		return strings.ToLower(quotes[i].EquipmentRequest.Equipment.Name) < strings.ToLower(quotes[j].EquipmentRequest.Equipment.Name)
	})
	return quotes
}

func SortQuotesByNameDesc(quotes []entities.Quote) []entities.Quote {
	sort.SliceStable(quotes, func(i, j int) bool {
		return strings.ToLower(quotes[i].EquipmentRequest.Equipment.Name) > strings.ToLower(quotes[j].EquipmentRequest.Equipment.Name)
	})
	return quotes
}

func SortQuotesByExpirationDateAsc(quotes []entities.Quote) []entities.Quote {
	sort.SliceStable(quotes, func(i, j int) bool {
		return quotes[i].EquipmentRequest.ExpirationDate < quotes[j].EquipmentRequest.ExpirationDate
	})
	return quotes
}

func SortQuotesByExpirationDateDesc(quotes []entities.Quote) []entities.Quote {
	sort.SliceStable(quotes, func(i, j int) bool {
		return quotes[i].EquipmentRequest.ExpirationDate > quotes[j].EquipmentRequest.ExpirationDate
	})
	return quotes
}

func UpdateListQuotesExpiryStatus(quotes []entities.Quote) []entities.Quote {
	for i := range quotes {
		currentTime := time.Now().Unix()
		if currentTime > quotes[i].EquipmentRequest.ExpirationDate {
			quotes[i].ExpiryStatus = "Expired"
		} else {
			quotes[i].ExpiryStatus = "Active"
		}
	}
	var activeQuotes []entities.Quote
	for _, quote := range quotes {
		if quote.ExpiryStatus == "Active" {
			activeQuotes = append(activeQuotes, quote)
		}
	}
	return activeQuotes
}

func ReturnLenOfQuotes(quotes []entities.Quote) int {
	return len(quotes)
}

func SortStatesWithoutRepeating(states []entities.State) []entities.State {
	sort.SliceStable(states, func(i, j int) bool {
		return states[i].Name < states[j].Name
	})
	var uniqueStates []entities.State
	for _, state := range states {
		if len(uniqueStates) == 0 {
			uniqueStates = append(uniqueStates, state)
			continue
		}
		if uniqueStates[len(uniqueStates)-1].Name != state.Name {
			uniqueStates = append(uniqueStates, state)
		}
	}
	states = uniqueStates
	return states
}

func SortCountiesWithoutRepeating(counties []entities.County) []entities.County {
	sort.SliceStable(counties, func(i, j int) bool {
		return counties[i].Name < counties[j].Name
	})
	var uniqueCounties []entities.County
	for i, county := range counties {
		if counties[i].Name == "" {
			continue
		}
		if len(uniqueCounties) == 0 {
			uniqueCounties = append(uniqueCounties, county)
			continue
		}
		if uniqueCounties[len(uniqueCounties)-1].Name != county.Name {
			uniqueCounties = append(uniqueCounties, county)
		}
	}
	counties = uniqueCounties
	return counties
}

func SortCitiesWithoutRepeating(cities []entities.City) []entities.City {
	sort.SliceStable(cities, func(i, j int) bool {
		return cities[i].Name < cities[j].Name
	})
	var uniqueCities []entities.City
	for _, city := range cities {
		if len(uniqueCities) == 0 {
			uniqueCities = append(uniqueCities, city)
			continue
		}
		if uniqueCities[len(uniqueCities)-1].Name != city.Name {
			uniqueCities = append(uniqueCities, city)
		}
	}
	cities = uniqueCities
	return cities
}

func ConvertToLowercase(counties []entities.County) []entities.County {
	for i := range counties {
		counties[i].Name = strings.ToLower(counties[i].Name)
	}
	return counties
}

func CapitalizeFirstLetterOfCounty(counties []entities.County) []entities.County {
	for i := range counties {
		counties[i].Name = strings.Title(counties[i].Name)
	}
	return counties
}
