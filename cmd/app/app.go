package app

import (
<<<<<<< HEAD:cmd/app/app.go
	"OnlineShopBackend/config"
=======
	"OnlineShopBackend/pkg/config"
	"OnlineShopBackend/pkg/filestorage"
	"OnlineShopBackend/pkg/logger"
>>>>>>> origin/Develop:pkg/app/app.go
	"context"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var GlobalApp *App

type Service interface {
	GetName() string
	Start(ctx context.Context) error
	ShutDown() error
}

type App struct {
	services []Service
	Log      *logger.Logger
	Fs       filestorage.FileStorager
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
<<<<<<< HEAD:cmd/app/app.go
		// TODO correct logger
		log.Printf("Error load config: %v. set default values", err)
	}
	type myString string
	var config myString = "config"
	ctx = context.WithValue(ctx, config, *cfg)
=======
		log.Println("Error load config. set default values")
	}
	a.Log = logger.NewLogger(cfg.LogLevel)

	ctx = context.WithValue(ctx, "config", *cfg)
>>>>>>> origin/Develop:pkg/app/app.go

	im := filestorage.NewOnDiskLocalStorage(cfg.DiskFileStoragePath)
	a.Fs = im

	for _, service := range a.services {
		service := service
		go func() {
			err := service.Start(ctx)
			if err != nil {
				a.Log.Logger.Panic("Error start service ", zap.Field{String: service.GetName()}, zap.Field{Interface: err})
			}
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c
	a.Log.Logger.Info("App shutdown...")

	for _, service := range a.services {
		service := service
		go func() {
			err := service.ShutDown()
			if err != nil {
				a.Log.Logger.Error("Error Shutdown service ", zap.Field{String: service.GetName()}, zap.Field{Interface: err})
			}
		}()
	}

	time.Sleep(time.Duration(cfg.Timeout) * time.Second)
}
