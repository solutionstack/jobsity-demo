package ws

import (
	"github.com/rs/zerolog"
	"github.com/solutionstack/lcache"
)

type WsService interface {
}

type service struct {
	logger zerolog.Logger
	cache  *lcache.Cache
}

func NewService(logger zerolog.Logger) WsService {
	return &service{
		cache:  lcache.NewCache(),
		logger: logger,
	}
}
