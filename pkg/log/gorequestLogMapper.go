package log

import (
	"github.com/parnurzeal/gorequest"
)

type gorequestLogMapper struct {
}

var loggerInstance *gorequestLogMapper

func (l *gorequestLogMapper) SetPrefix(prefix string) {
	//log.SetPrefix(prefix)

}
func (l *gorequestLogMapper) Printf(format string, v ...interface{}) {
	L.Infof(format, v)
}
func (l *gorequestLogMapper) Println(v ...interface{}) {
	L.Info(v)
}

func GetLogMapper() gorequest.Logger {

	if loggerInstance == nil {
		loggerInstance = &gorequestLogMapper{}
	}
	return loggerInstance
}
