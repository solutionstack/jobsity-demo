package main

import (
	logger "github.com/rs/zerolog"
	"github.com/solutionstack/jobsity-demo/cmd/server"
	handler "github.com/solutionstack/jobsity-demo/handlers/auth"
	wshandler "github.com/solutionstack/jobsity-demo/handlers/ws"
	svc "github.com/solutionstack/jobsity-demo/services/auth"
	wssvc "github.com/solutionstack/jobsity-demo/services/ws"
	"github.com/solutionstack/jobsity-demo/tickbot"
	"github.com/solutionstack/lcache"
	"os"
	"sync"
)

var Log = logger.New(os.Stdout)

func main() {
	var processManager sync.WaitGroup
	processManager.Add(3) //we'd wait for both server processes to complete

	fatalChan := make(chan error)

	//Dial into MQ broker
	mq := server.StartMQ()
	defer server.StopMQ(mq)
	mqChan, err := mq.Channel()
	if err != nil {
		panic(err)
	}

	//start stockbot worker
	stockBotMsgChan := make(chan []byte, 1)
	stockBotNew := tickbot.NewTickBot(Log, mqChan)
	go stockBotNew.StartWorker(&processManager, &stockBotMsgChan)

	cache := lcache.NewCache()
	clientPool := server.NewPool(Log, fatalChan, stockBotMsgChan) //client connection pool
	go clientPool.Start()                                         //start the client pool loop

	newService := svc.NewService(Log, cache)
	wsService := wssvc.NewService(Log, cache)

	newHandler := handler.NewHandler(Log, newService)
	wsHandler := wshandler.NewHandler(Log, stockBotNew, wsService)

	newRouter := server.NewRouter(newHandler)

	go server.StartWS(Log, wsHandler, &processManager, clientPool, fatalChan)
	go server.StartHTTPServer(newRouter, Log, &processManager, fatalChan)

	go func(fatalChan chan error) {
		for {
			select {
			case err := <-fatalChan:
				Log.Fatal().Err(err).Msg("Failed starting chat server")

			}
		}
	}(fatalChan)

	processManager.Wait()

}
