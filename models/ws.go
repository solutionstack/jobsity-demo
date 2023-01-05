package models

import "time"

type chatCommands string

const (
	MESSAGE      chatCommands = "MESSAGE"
	HISTORY      chatCommands = "HISTORY"
	STOCK_TICKER chatCommands = "STOCK_TICKER"
	ROOM_CREATE  chatCommands = "ROOM_CREATE"
	ROOM_READ    chatCommands = "ROOM_READ"
	USERS        chatCommands = "USERS"
	BAD_SESSION  chatCommands = "BAD_SESSION"
)

type WsMessage struct {
	Data       string       `json:"data,omitempty"`
	Command    chatCommands `json:"command"`
	Timestamp  time.Time    `json:"timestamp,omitempty"`
	Room       string       `json:"room,omitempty"`
	SessionKey string       `json:"sessionKey"`
	StockCode  string       `json:"stockCode"`
}

type WsResponse struct {
}
type WsUsers struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	SessionKey string `json:"sessionKey"`
}

type WsRoom struct {
	Name string `json:"name"`
}
type WsRoomMessage struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
