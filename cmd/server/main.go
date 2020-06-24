package main

import (
	"fmt"
	"github.com/alperhankendi/go-rest-api/internal/config"
	"github.com/labstack/echo/v4"
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

	var err error
	AppConfig, err = config.GetAllValues("./config/", fmt.Sprintf("config.%s", environment))

	if err != nil {
		panic(err)
	}
	instance = echo.New()
}
