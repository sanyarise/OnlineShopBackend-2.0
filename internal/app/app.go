package app

import (
	"OnlineShopBackend/config"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Service interface {
	GetName() string
	Start(ctx context.Context) error
	ShutDown() error
}

type App struct {
	logger   *zap.Logger
	services []Service
}

func NewApp(l *zap.Logger, serviceList []Service) *App {
	l.Debug("Enter in NewApp()")
	var s []Service
	s = append(s, serviceList...)
	return &App{logger: l, services: s}
}

func (app *App) Start() {
	app.logger.Debug("Enter in app Start()")
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		app.logger.Error(fmt.Sprintf("error load config: %v. set default values", err))
	}
	type myString string
	var config myString = "config"
	ctx = context.WithValue(ctx, config, *cfg)

	for _, service := range app.services {
		service := service
		go func() {
			err := service.Start(ctx)
			if err != nil {
				app.logger.Panic(fmt.Sprintf("error start service %s, %v", service.GetName(), err))
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c
	app.logger.Info("App shutdown...")

	for _, service := range app.services {
		service := service
		go func() {
			err := service.ShutDown()
			if err != nil {
				app.logger.Error(fmt.Sprintf("Error Shutdown service: %s, %v ", service.GetName(), err))
			}
		}()
	}

	time.Sleep(time.Duration(cfg.Timeout) * time.Second)
}
