package app

import (
	"OnlineShopBackend/config"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Service interface {
	GetName() string
	Start(ctx context.Context) error
	ShutDown() error
}

type App struct {
	services []Service
}

func NewApp(serviceList []Service) *App {
	var s []Service
	s = append(s, serviceList...)
	return &App{services: s}
}

func (a *App) Start() {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		// TODO correct logger
		log.Println("Error load config. set default values")
	}
	ctx = context.WithValue(ctx, "config", *cfg)

	for _, service := range a.services {
		service := service
		go func() {
			err := service.Start(ctx)
			if err != nil {
				// TODO correct logger
				log.Panicln("Error start service ", service.GetName(), err)
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c
	// TODO correct logger
	log.Println("App shutdown...")

	for _, service := range a.services {
		service := service
		go func() {
			err := service.ShutDown()
			if err != nil {
				// TODO correct logger
				log.Println("Error Shutdown service ", service.GetName(), err)
			}
		}()
	}

	time.Sleep(time.Duration(cfg.ShutDownTimeout) * time.Second)
}
