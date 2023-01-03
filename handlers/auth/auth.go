package auth

import (
	"github.com/rs/zerolog"
	"github.com/solutionstack/jobsity-demo/services/auth"
	"net/http"
)

type AuthHandler struct {
	logger zerolog.Logger
	svc    auth.Service
}

func NewHandler(logger zerolog.Logger, svc auth.Service) *AuthHandler {
	return &AuthHandler{
		logger: logger,
		svc:    svc,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}
