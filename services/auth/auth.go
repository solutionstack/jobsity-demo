package auth

import (
	"github.com/rs/zerolog"
	"github.com/solutionstack/lcache"
)

type Service interface {
}

type service struct {
	logger zerolog.Logger
	cache  *lcache.Cache
}

func NewService(logger zerolog.Logger) Service {
	return service{
		cache:  lcache.NewCache(),
		logger: logger,
	}
}
