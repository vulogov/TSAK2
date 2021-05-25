package tsak

import (
	"github.com/pieterclaerhout/go-log"

	tlog "github.com/vulogov/TSAK2/internal/log"
	"github.com/vulogov/TSAK2/internal/signal"
)

func Init() {
	tlog.Init()
	log.Debug("[ tsak2 ] tsak.Init() is reached")
	signal.InitSignal()
}
