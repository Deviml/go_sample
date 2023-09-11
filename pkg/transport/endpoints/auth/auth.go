package auth

import "github.com/go-kit/kit/endpoint"

type Services struct {
	Login LoginService
}

type Endpoints struct {
	Login endpoint.Endpoint
}

func NewAuth(svc Services) *Endpoints {
	return &Endpoints{Login: MakeLoginEndpoint(svc.Login)}
}
