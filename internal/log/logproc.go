package log

import (
	"fmt"
	"sync"
	"time"

	slog "github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/signal"
)

var N = 1000000
var Log = make(chan string, N)
var LOG_EVERY = (1 * time.Second)

func InitLogProc() {
	signal.Reserve(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		logproc()
		slog.Debug("Log process terminated")
	}(signal.WG())
}

func logproc() {
	var d string
	slog.Debug("Internal logging consumer started")
	for !signal.ExitRequested() {
		time.Sleep(LOG_EVERY)
		for len(Log) > 0 {
			d = <-Log
			fmt.Println(d)
		}
	}
}
