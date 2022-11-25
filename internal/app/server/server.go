package server

import (
	//sw "OnlineShopBackend/cmd/app"
	"OnlineShopBackend/internal/app/router"
	"context"
	"fmt"

	"go.uber.org/zap"
	//"reflect"
)

type HttpServer struct {
	ctx    context.Context
	port   string
	router *router.Router
	logger *zap.Logger
}

func NewServer(ctx context.Context, port string, router *router.Router, logger *zap.Logger) *HttpServer {
	logger.Debug("Enter in NewServer()")
	return &HttpServer{ctx: ctx, port: port, router: router, logger: logger}
}

func (server *HttpServer) GetName() string {
	server.logger.Debug("Enter in server GetName()")
	return "http server"
}

func (server *HttpServer) Start(ctx context.Context) error {
	server.logger.Debug("Enter in server Start()")
	server.ctx = ctx
	fmt.Println(server.port)

	err := server.router.Run(server.port)

	return err
}

func (server *HttpServer) ShutDown() error {
	server.logger.Debug("Enter in server ShubDown()")
	//TODO implement me
	panic("implement me")
}
