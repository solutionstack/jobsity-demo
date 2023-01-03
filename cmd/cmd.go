package main

import (
	logger "github.com/rs/zerolog"
	"github.com/solutionstack/jobsity-demo/cmd/server"

	handler "github.com/solutionstack/jobsity-demo/handlers/auth"
	svc "github.com/solutionstack/jobsity-demo/services/auth"
	"os"
)

var Log = logger.New(os.Stdout)

func main() {

	newService := svc.NewService(Log)
	newHandler := handler.NewHandler(Log, newService)
	newRouter := server.NewRouter(newHandler)

	if err := server.StartHTTPServer(newRouter, Log); err != nil {
		Log.Fatal().Err(err).Msg("Failed starting chat web-server")
	}

}
