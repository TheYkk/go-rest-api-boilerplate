package gommon

import (
	"github.com/alperhankendi/go-rest-api/pkg/log"
	"github.com/labstack/echo/v4"
)

type Adapter struct {
	echo.Logger
}

// FromGommon creates a new `log.Logger` from the provided entry
func FromGommon(logger echo.Logger) log.Logger {

	logger.SetPrefix("[MicroService]")
	return &Adapter{logger}
}

// you may extend your logger
func (l *Adapter) WithField(key string, val interface{}) log.Logger {
	return FromGommon(l.Logger)
}
