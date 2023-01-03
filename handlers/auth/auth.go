package auth

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/solutionstack/jobsity-demo/models"
	"github.com/solutionstack/jobsity-demo/services/auth"
	"github.com/solutionstack/jobsity-demo/utils"
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
	var body models.Signup

	if err := utils.UnmarshalRequestBody(r, &body); err != nil {
		a.logger.Error().Err(err).Msg(" bad request")
		utils.HandlerError(w, err, http.StatusBadRequest)
		return
	}

	if err := body.Validate(); err != nil {
		a.logger.Error().Err(err).Msg(" bad request")
		utils.HandlerError(w, err, http.StatusBadRequest)
		return
	}

	userID, err := a.svc.CreateUser(body)
	if err != nil {
		a.logger.Error().Err(err).Msg("service error")
		utils.HandlerError(w, errors.Wrap(err, "error while creating user"), http.StatusInternalServerError)
		return
	}

	res := models.RegisterResponse{
		ID: userID.String(),
	}
	resJson, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(resJson)
	return
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body models.Login

	if err := utils.UnmarshalRequestBody(r, &body); err != nil {
		utils.HandlerError(w, err, http.StatusBadRequest)
		return
	}

	if err := body.Validate(); err != nil {
		utils.HandlerError(w, err, http.StatusBadRequest)
		return
	}

	user, err := a.svc.ValidateLogin(body)
	if err != nil {
		utils.HandlerError(w, errors.Wrap(err, "error on login attempt"), http.StatusInternalServerError)
		return
	}
	if user == nil {
		utils.HandlerError(w, errors.New("no user record found"), http.StatusNotFound)
		return
	}

	resJson, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resJson)
}
