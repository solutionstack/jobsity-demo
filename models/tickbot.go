package models

type BotMQMessage struct {
	Room string `json:"room,omitempty"`
	Data string `json:"data"`
}
