package server

import (
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws/wsutil"
	"github.com/rs/zerolog"
	socketHandler "github.com/solutionstack/jobsity-demo/handlers/ws"
	"github.com/solutionstack/jobsity-demo/models"
	"net"
	"strings"
)

type Client struct {
	ID      string
	Conn    *net.Conn
	Pool    *Pool
	ErrChan chan<- error
	Handler *socketHandler.WsHandler
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		(*c.Conn).Close()
	}()

	for {
		msg, _, err := wsutil.ReadClientData(*c.Conn)
		if err != nil {
			if strings.Contains(err.Error(), "1001") { //peer closed
				return
			}
			c.ErrChan <- err
			return
		}

		//send new messages to command handler
		response, err := c.Handler.CommandHandler(msg)
		if err != nil {
			c.ErrChan <- err
			return
		}
		var respMsg models.WsMessage
		if err := json.Unmarshal(response, &respMsg); err != nil {
			c.ErrChan <- err
			return
		}

		//send message back to a single client or broadcast
		switch respMsg.Command {
		case models.ROOM_READ, models.BAD_SESSION, models.HISTORY:
			c.Pool.Logger.Info().Msg(fmt.Sprintf("Sending message to client. ID:%s", c.ID))
			err := wsutil.WriteServerMessage(*c.Conn, 0x1, response)
			if err != nil {
				c.Pool.ErrChan <- err
				return
			}
		default:
			c.Pool.Broadcast <- response
		}

	}
}

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Logger     zerolog.Logger
	ErrChan    chan<- error
}

func NewPool(logger zerolog.Logger, errChan chan<- error) *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Logger:     logger,
		ErrChan:    errChan,
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			pool.Logger.Info().Msg(fmt.Sprintf("new client joined. ID:%s ", client.ID))
			pool.Logger.Info().Msg(fmt.Sprintf("client pool size: %d", len(pool.Clients)))

			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			pool.Logger.Info().Msg(fmt.Sprintf("client disconnected. ID:%s ", client.ID))
			pool.Logger.Info().Msg(fmt.Sprintf("client pool size: %d", len(pool.Clients)))

			break
		case message := <-pool.Broadcast:
			pool.Logger.Info().Msg("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				err := wsutil.WriteServerMessage(*client.Conn, 0x1, message)

				if err != nil {
					pool.ErrChan <- err
					return
				}
			}
		}
	}
}
