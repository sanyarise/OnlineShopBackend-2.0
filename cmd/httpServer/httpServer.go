package httpServer

import (
	sw "OnlineShopBackend/cmd/app"
	"context"
	"fmt"
	"log"

	"go.uber.org/zap"
	//"reflect"
)

type HttpServer struct {
	ctx    context.Context
	port string
	router *sw.Router
	l *zap.Logger
}

func NewServer(ctx context.Context, port string, router *sw.Router, logger *zap.Logger) *HttpServer {
	log.Println("Enter in NewServer()")
	return &HttpServer{ctx: ctx,port: port, router: router, l: logger}
}

func (h *HttpServer) GetName() string {
	return "http server"
}

func (h *HttpServer) Start(ctx context.Context) error {
	h.ctx = ctx

	//cfg := h.ctx.Value("config")
	//port := reflect.ValueOf(cfg).FieldByName("Port").String()
	fmt.Println(h.port)

	err := h.router.Run(h.port)

	return err
}

func (h *HttpServer) ShutDown() error {
	//TODO implement me
	panic("implement me")
}
