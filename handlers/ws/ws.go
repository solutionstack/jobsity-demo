package ws

import (
	"github.com/rs/zerolog"
	"github.com/solutionstack/jobsity-demo/services/ws"
)

type WsHandler struct {
	logger zerolog.Logger
	svc    ws.WsService
}

func NewHandler(logger zerolog.Logger, svc ws.WsService) *WsHandler {
	return &WsHandler{
		logger: logger,
		svc:    svc,
	}
}

func (h *WsHandler) DefaultHandler(data []byte) ([]byte, error) {

	return data, nil
}
