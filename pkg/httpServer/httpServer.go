package httpServer

import (
	sw "OnlineShopBackend/cmd/onlineShopBackend/api"
	"context"
)

type HttpServer struct {
	ctx context.Context
}

func New() *HttpServer {
	return &HttpServer{}
}

func (h *HttpServer) GetName() string {
	return "http server"
}

func (h *HttpServer) Start(ctx context.Context) error {
	h.ctx = ctx

	router := sw.NewRouter()

	err := router.Run(":8000")

	return err
}

func (h *HttpServer) ShutDown() error {
	//TODO implement me
	panic("implement me")
}
