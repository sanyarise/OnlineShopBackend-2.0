package main

import (
	"OnlineShopBackend/pkg/app"
	"OnlineShopBackend/pkg/httpServer"
	"log"
	"os"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			os.Exit(1)
		}
	}()

	h := httpServer.New()
	var services []app.Service
	services = append(services, h)
	app.GlobalApp = app.NewApp(services)
	log.Printf("Server started")
	app.GlobalApp.Start()
}
