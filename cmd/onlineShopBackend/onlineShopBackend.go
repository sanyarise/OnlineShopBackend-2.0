package main

import (
	"OnlineShopBackend/config"
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
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Println("Start load configuration...")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("can't initialize configuration")
	}
	logger := logger.NewLogger(cfg.LogLevel)
	lsug := logger.Logger.Sugar()
	l := logger.Logger
	defer l.Sync()
	defer lsug.Sync()

	l.Info("Configuration sucessfully load")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	pgstore, err := repository.NewPgxStorage(ctx, lsug, cfg.DSN)
	if err != nil {
		log.Fatalf("can't initalize storage: %v", err)
	}
	itemStore := repository.NewItemRepo(pgstore, lsug)
	categoryStore := repository.NewCategoryRepo(pgstore, lsug)
	redis, err := cash.NewRedisCash(cfg.CashHost, cfg.CashPort, time.Duration(cfg.CashTTL), l)
	if err != nil {
		log.Fatalf("can't initialize cash: %v", err)
	}
	itemsCash := cash.NewItemsCash(redis, l)
	categoriesCash := cash.NewCategoriesCash(redis, l)

	itemUsecase := usecase.NewItemUsecase(itemStore, itemsCash, l)
	categoryUsecase := usecase.NewCategoryUsecase(categoryStore, categoriesCash, l)

	itemHandlers := handlers.NewItemHandlers(itemUsecase, l)
	categoryHandlers := handlers.NewCategoryHandlers(categoryUsecase, l)

	filestorage := filestorage.NewOnDiskLocalStorage(cfg.ServerURL, cfg.FsPath, l)
	delivery := delivery.NewDelivery(itemHandlers, categoryHandlers, l, filestorage)
	router := router.NewRouter(delivery, l)
	serverOptions := map[string]int{
		"ReadTimeout":       cfg.ReadTimeout,
		"WriteTimeout":      cfg.WriteTimeout,
		"ReadHeaderTimeout": cfg.ReadHeaderTimeout,
	}
	server := server.NewServer(cfg.Port, router, l, serverOptions)

	l.Debug("Start create cash...")
	_, err = categoryUsecase.GetCategoryList(ctx)
	if err != nil {
		l.Sugar().Errorf("error on create category cash: %w", err)
	}
	l.Info("Category list cash create success")

	_, err = itemUsecase.ItemsList(ctx, 0, 0)
	if err != nil {
		l.Sugar().Errorf("error on create items list cash: %w", err)
	}
	l.Info("Items list cash create success")

	server.Start()
	l.Info(fmt.Sprintf("Server start successful on port: %v", cfg.Port))

	<-ctx.Done()

	pgstore.ShutDown(cfg.Timeout)
	l.Info("Database connection stopped sucessful")

	redis.ShutDown(cfg.Timeout)
	l.Info("Cash connection stopped successful")

	server.ShutDown(cfg.Timeout)
	l.Info("Server stopped successful")
	cancel()
}
