package ws

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/solutionstack/jobsity-demo/models"
	"github.com/solutionstack/jobsity-demo/services/ws"
	"time"
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

func (h *WsHandler) CommandHandler(data []byte) ([]byte, error) {
	var wsMessage models.WsMessage

	if err := json.Unmarshal(data, &wsMessage); err != nil {
		return []byte(""), err
	}

	//is user session valud
	if !h.validateSession(wsMessage) {
		respMessage := models.WsMessage{
			Data:      "",
			Command:   models.BAD_SESSION,
			Timestamp: time.Now(),
			Room:      "",
		}
		jsonData, _ := json.Marshal(respMessage)

		return jsonData, nil
	}

	switch wsMessage.Command {
	case models.USERS: //fetch online users
		userList, _ := h.svc.GetUserList()
		userListJson, _ := json.Marshal(userList)

		respMessage := models.WsMessage{
			Data:      string(userListJson),
			Command:   models.USERS,
			Timestamp: time.Now(),
			Room:      "",
		}
		jsonData, _ := json.Marshal(respMessage)

		return jsonData, nil

	case models.ROOM_READ: //fetch room list
		roomList, _ := h.svc.GetRoomsList()
		roomListJson, _ := json.Marshal(roomList)

		respMessage := models.WsMessage{
			Data:      string(roomListJson),
			Command:   models.ROOM_READ,
			Timestamp: time.Now(),
			Room:      "",
		}
		jsonData, _ := json.Marshal(respMessage)

		return jsonData, nil

	case models.HISTORY: //fetch room messages
		msgs, _ := h.svc.GetMessages(wsMessage.Room)
		msgsJson, _ := json.Marshal(msgs)

		respMessage := models.WsMessage{
			Data:      string(msgsJson),
			Command:   models.HISTORY,
			Timestamp: time.Now(),
			Room:      "",
		}
		jsonData, _ := json.Marshal(respMessage)

		return jsonData, nil

	case models.MESSAGE: //fetch room messages
		msgs, err := h.svc.InsertMessage(wsMessage)
		if err != nil {
			fmt.Println(err)
			return []byte(""), nil
		}
		msgsJson, _ := json.Marshal(msgs)

		respMessage := models.WsMessage{
			Data:      string(msgsJson),
			Command:   models.MESSAGE,
			Timestamp: time.Now(),
			Room:      wsMessage.Room,
		}
		jsonData, _ := json.Marshal(respMessage)

		return jsonData, nil
	default:
		return []byte(""), nil

	}

}

func (h *WsHandler) validateSession(wsMessage models.WsMessage) bool {

	return h.svc.ValidateUserSession(wsMessage)
}

func (h *WsHandler) DefaultSetup() error {

	return h.svc.SetupDefaultConfig()
}
