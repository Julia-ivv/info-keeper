// Пакет logger реализует логер.
package logger

import (
	"go.uber.org/zap"
)

// ZapSugar предоставляет доступ к логеру.
var ZapSugar *zap.SugaredLogger

// NewLogger создает новый объект логера.
func NewLogger() *zap.SugaredLogger {
	log, errLog := zap.NewDevelopment()
	if errLog != nil {
		panic(errLog)
	}
	defer log.Sync()

	zapSugar := log.Sugar()

	return zapSugar
}
