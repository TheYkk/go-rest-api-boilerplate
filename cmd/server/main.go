package main

import (
	"fmt"
	"github.com/alperhankendi/go-rest-api/internal/basket"
	"github.com/alperhankendi/go-rest-api/internal/config"
	"github.com/alperhankendi/go-rest-api/pkg/graceful"
	"github.com/alperhankendi/go-rest-api/pkg/log"
	"github.com/alperhankendi/go-rest-api/pkg/mongoHelper"
	"github.com/labstack/echo/v4"
	"time"
)

var (
	port        = "5555"
	environment = "dev"
	debugMode   = false
	showHelp    = false
	instance    *echo.Echo

	AppConfig *config.Configuration
)

func main() {


	AppConfig, err := config.GetAllValues("./config/", fmt.Sprintf("config.%s", environment))

	if err != nil {
		log.L.Fatalf("Failed to parsing config file,Error:%s",err.Error())
	}

	instance = echo.New()
	instance.Debug = debugMode
	instance.HideBanner = !debugMode
	instance.HidePort = !debugMode
	log.L = instance.Logger


	// connect to the database
	 db, err := mongoHelper.ConnectDb(&AppConfig.MongoSettings)
	if err != nil {
		log.L.Fatalf("Database connection problem,Error:%v",err)
	}
	//bootstrapper for internal modules
	basket.RegisterHandlers(instance,basket.NewRepository(db))
	runAsService()
}

func runAsService() {

	log.L.Info("Service is starting with ", environment, " environment ", " and will serve on ", port, " port")

	go func() {
		if err := instance.Start(fmt.Sprintf(":%s", port)); err != nil {
			log.L.Fatalf("shutting down the server", err)
		}
	}()
	graceful.Shutdown(instance, 2*time.Second)
}
