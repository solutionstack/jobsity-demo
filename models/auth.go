package models

import (
	"github.com/pkg/errors"
	"net/mail"
	"strings"
)

type Signup struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID string `json:"id"'`
}

func (s *Signup) Validate() error {
	if strings.TrimSpace(s.Email) == "" {
		return errors.New("email field is required in the request body")
	}
	if strings.TrimSpace(s.Password) == "" {
		return errors.New("password field is required in the request body")
	}
	if strings.TrimSpace(s.FirstName) == "" {
		return errors.New("first-name field is required in the request body")
	}
	if _, err := mail.ParseAddress(s.Email); err != nil {
		return err
	}
	return nil

}
func (s *Login) Validate() error {
	if strings.TrimSpace(s.Email) == "" {
		return errors.New("email field is required in the request body")
	}
	if strings.TrimSpace(s.Password) == "" {
		return errors.New("password field is required in the request body")
	}
	return nil
}
