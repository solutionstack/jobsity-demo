package main

import (
	logger "github.com/rs/zerolog"
	"github.com/solutionstack/jobsity-demo/cmd/server"
	handler "github.com/solutionstack/jobsity-demo/handlers/auth"
	wshandler "github.com/solutionstack/jobsity-demo/handlers/ws"
	svc "github.com/solutionstack/jobsity-demo/services/auth"
	wssvc "github.com/solutionstack/jobsity-demo/services/ws"
	"os"
	"sync"
)

var Log = logger.New(os.Stdout)

func main() {
	var processManager sync.WaitGroup
	processManager.Add(2) //wait for both server processes to complete

	newService := svc.NewService(Log)
	wsService := wssvc.NewService(Log)
	newHandler := handler.NewHandler(Log, newService)
	wsHandler := wshandler.NewHandler(Log, wsService)
	newRouter := server.NewRouter(newHandler)

	fatalChan := make(chan error)
	go server.StartWS(Log, wsHandler, &processManager, fatalChan)
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
