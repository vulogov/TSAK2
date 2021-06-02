package dsl

import (
	"time"

	"github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/signal"
)

func PanicOn(err error) {
	if err != nil {
		signal.ExitRequest()
		log.Debug("We detected a fatal condition in TSAK2 application. Sending ExitRequest()")
		time.Sleep(time.Second)
		log.Fatalf("Fatal condition: %v", err)
	}
}
