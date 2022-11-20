package httpServer

import (
	sw "OnlineShopBackend/cmd/app"
	"context"
	"fmt"
	"log"
	//"reflect"
)

type HttpServer struct {
	ctx    context.Context
	router *sw.Router
}

func NewServer(ctx context.Context, router *sw.Router) *HttpServer {
	log.Println("Enter in NewServer()")
	return &HttpServer{ctx: ctx, router: router}
}

func (h *HttpServer) GetName() string {
	return "http server"
}

func (h *HttpServer) Start(ctx context.Context, port string) error {
	h.ctx = ctx

	//cfg := h.ctx.Value("config")
	//port := reflect.ValueOf(cfg).FieldByName("Port").String()
	fmt.Println(port)

	err := h.router.Run(port)

	return err
}

func (h *HttpServer) ShutDown() error {
	//TODO implement me
	panic("implement me")
}
