package ws

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/solutionstack/jobsity-demo/models"
	"github.com/solutionstack/lcache"
	"time"
)

const (
	defaultWsRoom      = "Public"
	defaultWsUser      = "TickBot"
	defaultWsUserEmail = "tickbot@bot.net"
	wsUsersCacheKey    = "WS_USER_DATA"
	wsRoomsCacheKey    = "WS_ROOM_DATA"
	SessionPrefix      = "Session_"
	maxRomMessageCount = 50
)

type WsService interface {
	SetupDefaultConfig() error
	GetUserList() ([]models.WsUsers, error)
	ValidateUserSession(msg models.WsMessage) bool
	GetRoomsList() ([]models.WsRoom, error)
	GetMessages(roomName string) ([]models.WsRoomMessage, error)
	InsertMessage(roomName models.WsMessage) (models.WsRoomMessage, error)
}

type service struct {
	logger zerolog.Logger
	cache  *lcache.Cache
}

func NewService(logger zerolog.Logger, c *lcache.Cache) WsService {
	return &service{
		cache:  c,
		logger: logger,
	}
}

// SetupDefaultConfig setup default users and rooms
func (s *service) SetupDefaultConfig() error {
	var wsUsers []models.WsUsers
	var wsRooms []models.WsRoom
	var wsRoomMessages []models.WsRoomMessage

	wsUsers = append(wsUsers, models.WsUsers{
		Name:  defaultWsUser,
		Email: defaultWsUserEmail,
	})
	data, _ := json.Marshal(wsUsers)
	s.cache.Write(wsUsersCacheKey, string(data))

	wsRooms = append(wsRooms, models.WsRoom{
		Name: defaultWsRoom,
	})
	data, _ = json.Marshal(wsRooms)
	s.cache.Write(wsRoomsCacheKey, string(data))

	wsRoomMessages = append(wsRoomMessages, models.WsRoomMessage{
		Name:  defaultWsUser,
		Email: defaultWsUserEmail,
		Data: "<b>Hello</b>, remember you can use the \"Stock Ticker\" bot by sending a slash command in the format<br/>" +
			"<b>/stock=stock_code</b",
		Timestamp: time.Now(),
	})
	data, _ = json.Marshal(wsRoomMessages)
	for _, r := range wsRooms {
		s.cache.Write(r.Name+"_Messages", string(data))

	}

	return nil
}

func (s *service) GetUserList() ([]models.WsUsers, error) {
	var wsUsers []models.WsUsers
	result := s.cache.Read(wsUsersCacheKey)

	if result.Value != nil {
		json.Unmarshal([]byte(result.Value.(string)), &wsUsers)
	}

	return wsUsers, nil
}

func (s *service) GetRoomsList() ([]models.WsRoom, error) {
	var rooms []models.WsRoom
	result := s.cache.Read(wsRoomsCacheKey)

	if result.Value != nil {
		json.Unmarshal([]byte(result.Value.(string)), &rooms)
	}

	return rooms, nil
}

func (s *service) GetMessages(roomName string) ([]models.WsRoomMessage, error) {
	var rooms []models.WsRoomMessage
	result := s.cache.Read(roomName + "_Messages")

	if result.Value != nil {
		json.Unmarshal([]byte(result.Value.(string)), &rooms)
	}

	return rooms, nil
}

func (s *service) InsertMessage(wsMsg models.WsMessage) (models.WsRoomMessage, error) {
	currentRoomMessages, err := s.GetMessages(wsMsg.Room)
	if err != nil {
		return models.WsRoomMessage{}, err
	}

	var newMsg models.WsRoomMessage
	err = json.Unmarshal([]byte(wsMsg.Data), &newMsg)
	if err != nil {
		return models.WsRoomMessage{}, err
	}

	currentRoomMessages = append(currentRoomMessages, newMsg)

	var data []byte
	if len(currentRoomMessages) >= maxRomMessageCount {
		data, _ = json.Marshal(currentRoomMessages[0:maxRomMessageCount])

	} else {
		data, _ = json.Marshal(currentRoomMessages)

	}
	s.cache.Write(wsMsg.Room+"_Messages", string(data))

	return newMsg, nil

}
func (s *service) ValidateUserSession(msg models.WsMessage) bool {
	result := s.cache.Read(SessionPrefix + msg.SessionKey)
	return result.Error == nil
}
