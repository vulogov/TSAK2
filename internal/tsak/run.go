package tsak

import (
	"fmt"
	"os"
	"time"

	"github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/conf"
	"github.com/vulogov/TSAK2/internal/dsl"
	"github.com/vulogov/TSAK2/internal/signal"
)

func Run() {
	Init()
	log.Debug("[ tsak2 ] tsak.Run() is reached")
	cfg := dsl.InitDSL()
	env := dsl.MakeEnvironment(cfg)
	if *conf.BootStrap != "" {
		log.Debugf("TSAK_script bootstrap file: %v", *conf.BootStrap)
		file, err := os.Open(*conf.BootStrap)
		err = env.LoadFile(file)
		dsl.PanicOn(err)
		_, err = env.Run()
		dsl.PanicOn(err)
	}
	for _, fn := range *conf.Scripts {
		log.Debugf("Loading script from %s in TSAK environment", fn)
		file, err := os.Open(fn)
		err = env.LoadFile(file)
		dsl.PanicOn(err)
	}
	log.Debugf("TSAK-script environment is running")
	res, err := env.Run()
	dsl.PanicOn(err)
	if *conf.ShowResult {
		log.Debugf("Displaying result were requested")
		fmt.Println(res.SexpString(nil))
	} else {
		log.Debugf("Displaying result was not requested")
	}
	if *conf.ERloop {
		log.Debug("ExitRequest event loop reached")
		for !signal.ExitRequested() {
			time.Sleep(100 * time.Millisecond)
		}
	} else if *conf.WGloop {
		log.Debug("WorkGroup event loop reached")
		signal.Loop()
	} else {
		log.Debug("EvenLoop was not requested")
	}
	signal.ExitRequest()
}
