package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServer interface {
	Start() error
	ShutDown(ctx context.Context) error
}

type App struct {
	httpServer HttpServer
}

func NewApp(httpServer HttpServer) *App {
	return &App{httpServer: httpServer}
}

func (a *App) Start() {
	go func() {
		err := a.httpServer.Start()
		if err != nil {

			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c
	log.Println("App shutdown...")

	// TODO вынести в конфиг настройку таймаута
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.ShutDown(ctx)
}

func (a *App) ShutDown(ctx context.Context) {
	err := a.httpServer.ShutDown(ctx)
	if err != nil {
		log.Println("Error shutdown http server", err)
	}
}
