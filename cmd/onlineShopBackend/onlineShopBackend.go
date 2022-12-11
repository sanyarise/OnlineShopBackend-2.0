package main

import (
	"OnlineShopBackend/config"
	"OnlineShopBackend/internal/app"
	"OnlineShopBackend/internal/app/logger"
	"OnlineShopBackend/internal/app/router"
	"OnlineShopBackend/internal/app/server"
	"OnlineShopBackend/internal/delivery"
	"OnlineShopBackend/internal/filestorage"
	"OnlineShopBackend/internal/handlers"
	"OnlineShopBackend/internal/repository"
	"OnlineShopBackend/internal/repository/cash"
	"OnlineShopBackend/internal/usecase"
	"context"
	"log"
	"os"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			os.Exit(1)
		}
	}()
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("can't initialize configuration")
	}
	logger := logger.NewLogger(cfg.LogLevel)
	lsug := logger.Logger.Sugar()
	l := logger.Logger
	pgstore, err := repository.NewPgxStorage(ctx, lsug, cfg.DSN)
	if err != nil {
		log.Fatalf("can't initalize storage: %v", err)
	}
	itemStore := repository.NewItemRepo(pgstore, lsug)
	categoryStore := repository.NewCategoryRepo(pgstore, lsug)
	cash, err := cash.NewRedisCash(cfg.CashHost, cfg.CashPort, time.Duration(cfg.CashTTL), l)
	if err != nil {
		log.Fatalf("can't initialize cash: %v", err)
	}

	itemUsecase := usecase.NewItemUsecase(itemStore, cash, l)
	categoryUsecase := usecase.NewCategoryUsecase(categoryStore, l)

	itemHandlers := handlers.NewItemHandlers(itemUsecase, l)
	categoryHandlers := handlers.NewCategoryHandlers(categoryUsecase, l)

	filestorage := filestorage.NewOnDiskLocalStorage(cfg.ServerURL, cfg.FsPath, l)
	delivery := delivery.NewDelivery(itemHandlers, categoryHandlers, l, filestorage)
	router := router.NewRouter(delivery, l)
	server := server.NewServer(ctx, cfg.Port, router, l)
	err = server.Start(ctx)
	if err != nil {
		log.Fatalf("can't start server: %v", err)
	}
	var services []app.Service

	services = append(services, server)
	a := app.NewApp(l, services)
	log.Printf("Server started")
	a.Start()
}
