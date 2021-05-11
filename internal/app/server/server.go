package server

import (
	"context"
	"fmt"
	"go-starter-project/internal/app/config"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type server struct {
	server *http.Server
}

func NewServer(handler http.Handler) *server {
	svr := &server{
		server: &http.Server{
			Addr:              fmt.Sprint(":", config.Conf.Server.Port),
			Handler:           handler,
			ReadHeaderTimeout: config.Conf.Server.ReadHeaderTimeout,
			WriteTimeout:      config.Conf.Server.WriteTimeout,
		},
	}
	return svr
}

func (s *server) serveWithGratefulShutdown() {
	// Create listener for the 'SIGTERM' from kernel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Wait for 'SIGTERM' from kernel
	<-quit

	// Create the cancelable context for help cancel the halted shutdown process
	srvCtx, srvCancel := context.WithTimeout(context.Background(), config.Conf.Server.GracefulShutdownTime)
	defer srvCancel()

	// Perform shutdown then wait until the server finished the shutdown
	// process or the timeout had been reached
	log.Println("Shutting down HTTP server")
	if err := s.server.Shutdown(srvCtx); err != nil {
		log.Panic("HTTP server shutdown with error")
	}
}

func (s *server) ListenAndServe() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic("Listen & Serve Failed !!!")
		}
	}()
	s.serveWithGratefulShutdown()
}
