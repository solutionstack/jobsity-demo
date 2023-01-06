package tickbot

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/solutionstack/jobsity-demo/models"
	"github.com/streadway/amqp"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	mqQueueName = "tickbot"
)

type TickBot interface {
	StartWorker(pm *sync.WaitGroup, respChan *chan []byte)
	BotMessage(data models.WsMessage) error
}

type bot struct {
	httpClient   *http.Client
	logger       zerolog.Logger
	mqChannel    *amqp.Channel
	ShutdownChan chan os.Signal
}

func NewTickBot(logger zerolog.Logger, mqChan *amqp.Channel) TickBot {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	_, err := mqChan.QueueDeclare(
		mqQueueName, // queue name
		true,        // durable
		false,       // auto delete
		false,       // exclusive
		false,       // no wait
		nil,         // arguments
	)
	if err != nil {
		panic(err)
	}

	return &bot{
		logger:       logger,
		mqChannel:    mqChan,
		ShutdownChan: stop,
		httpClient: &http.Client{
			Timeout: time.Duration(15) * time.Second,
		},
	}
}

func (b *bot) BotMessage(data models.WsMessage) error {
	botMsg := models.BotMQMessage{
		Room: data.Room,
		Data: data.Data,
	}
	botMsgJson, err := json.Marshal(botMsg)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        botMsgJson,
	}

	// Attempt to publish a message to the queue.
	if err := b.mqChannel.Publish(
		"",          // exchange
		mqQueueName, // queue name
		false,       // mandatory
		false,       // immediate
		message,     // message to publish
	); err != nil {
		return err
	}

	return nil
}

func (b *bot) StartWorker(pm *sync.WaitGroup, respChan *chan []byte) {
	b.logger.Info().Msg("starting bot worker")

	messagesChan, err := b.mqChannel.Consume(
		mqQueueName, // queue name
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no local
		false,       // no wait
		nil,         // arguments
	)
	if err != nil {
		b.logger.Error().Err(err).Msg("error reading message from mq")
	}

	go func() {
		for message := range messagesChan {
			b.logger.Info().Msgf(" >BOT: new queue message: %s\n", message.Body)
			go b.processMessage(message.Body, respChan)
		}
	}()
	<-b.ShutdownChan
	pm.Done()
	b.logger.Info().Msg("stopping bot worker")
}

func (b *bot) processMessage(msg []byte, respChan *chan []byte) {
	errorState := false

	var botMsg models.BotMQMessage
	json.Unmarshal(msg, &botMsg)

	var msgData models.WsRoomMessage
	json.Unmarshal([]byte(botMsg.Data), &msgData)

	stockCode := strings.Replace(msgData.Data, "/stock=", "", 1)
	url := "https://stooq.com/q/l/?s=" + stockCode + "&f=sd2t2ohlcv&h&e=csv"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		b.logger.Error().Err(err).Msg("processMessage:NewRequest")
		errorState = true
	}

	resp, err := b.httpClient.Do(req)
	if err != nil {
		b.logger.Error().Err(err).Msg("processMessage:Do")
		errorState = true
	}

	//read csv response from stock quote provider
	r := csv.NewReader(resp.Body)
	rows, _ := r.ReadAll()
	if err != nil {
		b.logger.Error().Err(err).Msg("processMessage:ReadAll")
		errorState = true
	}

	if errorState {
		msgData.Data = "Could not fetch stock quotes"

	} else {
		if strings.Contains(rows[1][6], "N/D") {
			msgData.Data = "Unable to fetch quotes for " + strings.ToUpper(stockCode) + " at this time"
		} else {
			msgData.Data = strings.ToUpper(stockCode) + " quote is $" + rows[1][6] + " per share"

		}
	}
	msgDataJson, _ := json.Marshal(msgData)
	msgData.Timestamp = time.Now()

	wsMsg := models.WsMessage{
		Data:       string(msgDataJson),
		Command:    models.STOCK_TICKER,
		Timestamp:  time.Now(),
		Room:       botMsg.Room,
		SessionKey: "",
		StockCode:  rows[1][0],
	}

	wsMsgJson, _ := json.Marshal(wsMsg)
	fmt.Println(string(wsMsgJson))
	*respChan <- wsMsgJson
}
