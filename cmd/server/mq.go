package server

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

func StartMQ() *amqp.Connection {
	AMQ_HOST := os.Getenv("AMQ_HOST")
	AMQ_USER := os.Getenv("AMQ_USER")
	AMQ_PASS := os.Getenv("AMQ_PASS")

	connection, err := amqp.Dial("amqp://" + AMQ_USER + ":" + AMQ_PASS + "@" + AMQ_HOST + "/" + AMQ_USER)
	if err != nil {
		panic(err)
	}
	fmt.Println("rabbitmq connected")
	return connection

}
func StopMQ(mq *amqp.Connection) {
	err := mq.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("rabbitmq disconnected")
}
