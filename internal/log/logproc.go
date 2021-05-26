package log

import (
	"sync"
	"time"

	slog "github.com/pieterclaerhout/go-log"

	"github.com/Jeffail/gabs/v2"

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
	var msg string
	slog.Debug("Internal logging consumer started")
	for !signal.ExitRequested() {
		time.Sleep(LOG_EVERY)
		for len(Log) > 0 {
			msg = <-Log
			displayLog(msg)
		}
	}
	slog.Debug("Internal logging consumer terminated")
}

func displayLog(msg string) {
	jsonParsed, err := gabs.ParseJSON([]byte(msg))
	if err != nil {
		return
	}
	out, ok := jsonParsed.Path("out").Data().(string)
	if !ok {
		return
	}
	msgOut, ok := jsonParsed.Path("msg").Data().(string)
	if !ok {
		return
	}
	if out == "debug" {
		slog.Debug(msgOut)
	} else if out == "info" {
		slog.Info(msgOut)
	} else if out == "warning" {
		slog.Warn(msgOut)
	} else if out == "error" {
		slog.Error(msgOut)
	} else if out == "fatal" {
		slog.Fatal(msgOut)
	} else {
		slog.Debug(msgOut)
	}
}
