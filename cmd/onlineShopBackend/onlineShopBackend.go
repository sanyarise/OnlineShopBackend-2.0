package main

import (
	"OnlineShopBackend/config"
	"OnlineShopBackend/internal/app/logger"
	"OnlineShopBackend/internal/app/router"
	"OnlineShopBackend/internal/app/server"
	"OnlineShopBackend/internal/delivery"
	"OnlineShopBackend/internal/filestorage"
	"OnlineShopBackend/internal/repository"
	"OnlineShopBackend/internal/repository/cash"
	"OnlineShopBackend/internal/usecase"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
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

	l.Info("Configuration sucessfully load")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	pgstore, err := repository.NewPgxStorage(ctx, lsug, cfg.DSN)
	if err != nil {
		log.Fatalf("can't initalize storage: %v", err)
	}
	itemStore := repository.NewItemRepo(pgstore, lsug)
	categoryStore := repository.NewCategoryRepo(pgstore, lsug)

	cartStore := repository.NewCartStore(pgstore, lsug)

	redis, err := cash.NewRedisCash(cfg.CashHost, cfg.CashPort, time.Duration(cfg.CashTTL), l)
	if err != nil {
		log.Fatalf("can't initialize cash: %v", err)
	}
	itemsCash := cash.NewItemsCash(redis, l)
	categoriesCash := cash.NewCategoriesCash(redis, l)

	itemUsecase := usecase.NewItemUsecase(itemStore, itemsCash, l)
	categoryUsecase := usecase.NewCategoryUsecase(categoryStore, categoriesCash, l)

	cartUsecase := usecase.NewCartUseCase(cartStore, l)

	filestorage := filestorage.NewOnDiskLocalStorage(cfg.ServerURL, cfg.FsPath, l)
	delivery := delivery.NewDelivery(itemUsecase, categoryUsecase, cartUsecase, l, filestorage)

	router := router.NewRouter(delivery, l)
	serverOptions := map[string]int{
		"ReadTimeout":       cfg.ReadTimeout,
		"WriteTimeout":      cfg.WriteTimeout,
		"ReadHeaderTimeout": cfg.ReadHeaderTimeout,
	}
	server := server.NewServer(cfg.Port, router, l, serverOptions)

	err = createCashOnStartService(ctx, categoryUsecase, itemUsecase, l)
	if err != nil {
		l.Sugar().Fatalf("error on create cash on start: %v", err)
	}

	server.Start()
	l.Info(fmt.Sprintf("Server start successful on port: %v", cfg.Port))

	<-ctx.Done()

	err = pgstore.ShutDown(cfg.Timeout)
	if err != nil {
		l.Error(err.Error())
	} else {
		l.Info("Database connection stopped sucessful")
	}

	err = redis.ShutDown(cfg.Timeout)
	if err != nil {
		l.Error(err.Error())
	} else {
		l.Info("Cash connection stopped successful")
	}

	err = server.ShutDown(cfg.Timeout)
	if err != nil {
		l.Error(err.Error())
	} else {
		l.Info("Server stopped successful")
	}

	cancel()
}

func createCashOnStartService(ctx context.Context, categoryUsecase usecase.ICategoryUsecase, itemUsecase usecase.IItemUsecase, l *zap.Logger) error {
	l.Debug("Enter in main createCashOnStartService")
	l.Debug("Start create cash...")
	categoryList, err := categoryUsecase.GetCategoryList(ctx)
	if err != nil {
		l.Sugar().Errorf("error on create category cash: %w", err)
		return err
	}
	l.Info("Category list cash create success")

	limitOptions := map[string]int{"offset": 0, "limit": 0}
	listOptions := []map[string]string{
		{"sortType": "name", "sortOrder": "asc"},
		{"sortType": "name", "sortOrder": "desc"},
		{"sortType": "price", "sortOrder": "asc"},
		{"sortType": "price", "sortOrder": "desc"},
	}
	for _, sortOptions := range listOptions {
		_, err = itemUsecase.ItemsList(ctx, limitOptions, sortOptions)
		if err != nil {
			l.Sugar().Errorf("error on create items list cash: %w", err)
			return err
		}
		l.Info("Items list cash create success")

		for _, category := range categoryList {
			_, err := itemUsecase.GetItemsByCategory(ctx, category.Name, limitOptions, sortOptions)
			if err != nil {
				l.Sugar().Errorf("error on create items list in category: %s cash: %w", category.Name, err)
				return err
			}
		}
	}
	l.Info("Items lists in categories cash create success")
	return nil
}
