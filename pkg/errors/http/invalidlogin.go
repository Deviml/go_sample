package http

import (
	"net/http"
)

type InvalidLoginError struct {
}

type ExistingCredentials struct {
}

func (e ExistingCredentials) StatusCode() int {
	// Credentials already exist
	return http.StatusConflict
}

func (e ExistingCredentials) Error() string {
	return "Email already exist"
}

func (i InvalidLoginError) StatusCode() int {
	return http.StatusUnauthorized
}

func (i InvalidLoginError) Error() string {
	return "wrong credentials"
}
