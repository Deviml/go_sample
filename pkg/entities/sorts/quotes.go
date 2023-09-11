package sorts

import "github.com/Equiphunter-com/equipment-hunter-api/pkg/transport/requests"

type ListQuotes struct {
	Criteria  string
	Latitude  float32
	Longitude float32
}

func MakeListQuotesFromRequest(listRequest requests.ListQuotesRequest) ListQuotes {
	return ListQuotes{
		Criteria:  listRequest.Sort,
		Latitude:  listRequest.Latitude,
		Longitude: listRequest.Longitude,
	}
}
