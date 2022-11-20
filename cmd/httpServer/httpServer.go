package httpServer

import (
	sw "OnlineShopBackend/cmd/app"
	"context"
	"fmt"
	"reflect"
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

	cfg := h.ctx.Value("config")
	port := reflect.ValueOf(cfg).FieldByName("Port").String()
	fmt.Println(port)

	err := router.Run(port)

	return err
}

func (h *HttpServer) ShutDown() error {
	//TODO implement me
	panic("implement me")
}
