package server

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	srv    http.Server
	logger *zap.Logger
}

// NewServer returns new server with configured parameters
func NewServer(addr string, handler http.Handler, logger *zap.Logger, timeouts map[string]int) *Server {
	logger.Debug("Enter in NewServer()")
	server := &Server{}

	server.srv = http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadTimeout:       time.Duration(timeouts["ReadTimeout"]) * time.Second,
		WriteTimeout:      time.Duration(timeouts["WriteTimeout"]) * time.Second,
		ReadHeaderTimeout: time.Duration(timeouts["ReadHeaderTimeout"]) * time.Second,
	}
	server.logger = logger
	return server
}

// Start begin server work
func (server *Server) Start() {
	server.logger.Debug("Enter in server Start()")
	go server.srv.ListenAndServe()
}

// ShutDown stop the server
func (server *Server) ShutDown(ctx context.Context, timeout int) {
	server.logger.Debug("Enter in server ShubDown()")
	ctxWithTimiout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	server.srv.Shutdown(ctxWithTimiout)
	cancel()
}
