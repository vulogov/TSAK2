package log

import (
	"github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/conf"
)

func Init() {
	if *conf.Debug {
		log.DebugMode = true
	} else {
		log.DebugMode = false
	}
	log.PrintTimestamp = true
	log.TimeFormat = "15:04:05.000"
	if *conf.Color {
		log.PrintColors = true
	} else {
		log.PrintColors = false
	}
	log.Debug("[ TSAK2 ] console log is configured")
	InitLogProc()
}
