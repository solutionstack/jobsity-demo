package server

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/rs/zerolog"
	socketHandler "github.com/solutionstack/jobsity-demo/handlers/ws"
	"github.com/solutionstack/jobsity-demo/utils"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	wsAddress = "6001"
)

func StartWS(logger zerolog.Logger, handler *socketHandler.WsHandler, pm *sync.WaitGroup, pool *Pool, errChan chan<- error) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	mux := http.NewServeMux()

	srv := &http.Server{Addr: ":" + wsAddress, Handler: mux}
	go func() {
		logger.Log().
			Str("address", srv.Addr).
			Int("pid", os.Getpid()).
			Msg("chat socket server listening")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("listenAndServe failed: %v", err)
		}
	}()

	go func() {
		//CLEANUP
		<-stop

		//close server
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}

		logger.Log().Msg("chat socket server shutdown")
		pm.Done()
	}()

	//default socket handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			errChan <- err
			return
		}

		//add current client to client pool
		client := &Client{
			ID:      utils.RandStringBytesRmndr(6),
			Conn:    &conn,
			Pool:    pool,
			ErrChan: errChan,
			Handler: handler,
		}
		pool.Register <- client

		client.Read() //setup reading for current client

	})

	//default chat room settings
	err := handler.DefaultSetup()
	if err != nil {
		errChan <- err
		return
	}

}
