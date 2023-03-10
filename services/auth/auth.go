package auth

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/solutionstack/jobsity-demo/models"
	"github.com/solutionstack/lcache"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(user models.Signup) (uuid.UUID, error)
	ValidateLogin(login models.Login) (*models.UserRecord, error)
}

const (
	userDbPrefix     = "User_"
	SessionPrefix    = "Session_"
	passwordHashCost = 15
	wsUsersCacheKey  = "WS_USER_DATA"
	wsRoomsCacheKey  = "WS_ROOM_DATA"
)

type service struct {
	logger zerolog.Logger
	cache  *lcache.Cache
}

func NewService(logger zerolog.Logger, c *lcache.Cache) Service {
	return &service{
		cache:  c,
		logger: logger,
	}
}

func (s *service) CreateUser(user models.Signup) (uuid.UUID, error) {

	emailHash := md5.Sum([]byte(user.Email))
	passwdHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), passwordHashCost)
	if err != nil {
		return uuid.Nil, err
	}

	if exists := s.cache.Read(userDbPrefix + fmt.Sprintf("%x", emailHash)); exists.Value != nil {
		return uuid.Nil, errors.New("user already exist")
	}

	userRecord := models.UserRecord{
		ID: uuid.New(),
		Signup: models.Signup{
			FirstName: user.FirstName,
			Email:     user.Email,
			Password:  string(passwdHash),
		},
	}

	j, err := json.Marshal(userRecord)
	if err != nil {
		return uuid.Nil, err
	}

	s.cache.Write(userDbPrefix+fmt.Sprintf("%x", emailHash), string(j))

	return userRecord.ID, nil
}
func (s *service) ValidateLogin(login models.Login) (*models.UserRecord, error) {
	emailHash := md5.Sum([]byte(login.Email))
	userRecordKey := userDbPrefix + fmt.Sprintf("%x", emailHash)

	userRecord := s.cache.Read(userRecordKey)
	if userRecord.Error != nil && userRecord.Error == lcache.KeyNotFoundError {
		return nil, nil
	}

	var user models.UserRecord
	err := json.Unmarshal([]byte(userRecord.Value.(string)), &user)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return nil, errors.New("password mismatch")
	}

	sk := s.initSession(login)
	user.Password = "--redacted--"
	user.SessionKey = sk
	s.addUserToWsList(user)

	return &user, nil
}

func (s *service) initSession(login models.Login) string {
	emailHash := md5.Sum([]byte(login.Email))
	sessionKey := SessionPrefix + fmt.Sprintf("%x", emailHash)

	s.cache.Write(sessionKey, login.Email)
	return sessionKey
}

// add to websocket online users cache
func (s *service) addUserToWsList(user models.UserRecord) {
	var wsUsers []models.WsUsers
	result := s.cache.Read(wsUsersCacheKey)

	if result.Value != nil {
		json.Unmarshal([]byte(result.Value.(string)), &wsUsers)
	}

	wsUsers = append(wsUsers, models.WsUsers{
		Name:       user.FirstName,
		Email:      user.Email,
		SessionKey: user.SessionKey,
	})
	data, _ := json.Marshal(wsUsers)

	s.cache.Write(wsUsersCacheKey, string(data))
}
