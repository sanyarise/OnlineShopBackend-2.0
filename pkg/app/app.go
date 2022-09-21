package app

import (
	"OnlineShopBackend/pkg/config"
	"OnlineShopBackend/pkg/logger"
	"context"
	"go.uber.org/zap"
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
	log      *logger.Logger
}

func NewApp(serviceList []Service) *App {
	var s []Service
	s = append(s, serviceList...)
	return &App{services: s}
}

func (a *App) Start() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Println("Error load config. set default values")
	}
	a.log = logger.NewLogger(cfg.LogLevel)

	ctx = context.WithValue(ctx, "config", *cfg)

	for _, service := range a.services {
		service := service
		go func() {
			err := service.Start(ctx)
			if err != nil {
				a.log.Logger.Panic("Error start service ", zap.Field{String: service.GetName()}, zap.Field{Interface: err})
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c
	a.log.Logger.Info("App shutdown...")

	for _, service := range a.services {
		service := service
		go func() {
			err := service.ShutDown()
			if err != nil {
				a.log.Logger.Error("Error Shutdown service ", zap.Field{String: service.GetName()}, zap.Field{Interface: err})
			}
		}()
	}

	time.Sleep(time.Duration(cfg.Timeout) * time.Second)
}
