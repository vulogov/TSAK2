package dsl

import (
	"github.com/glycerine/zygomys/zygo"
	"github.com/pieterclaerhout/go-log"
)

func TsakPipeSetup(cfg *zygo.ZlispConfig, env *zygo.Zlisp) {
	dPkg := `(def pipe (package "pipe"
     { Create := pipecreate;
       Send := pipesend;
       Len := pipelen;
       Recv := piperecv;
     }
  ))`
	_, err := env.EvalString(dPkg)
	PanicOn(err)
}

func TsakStandardSetup(cfg *zygo.ZlispConfig, env *zygo.Zlisp) {
	log.Debug("DSL setup for TSAK is reached")
	TsakLogSetup(cfg, env)
	TsakPipeSetup(cfg, env)
	TelemetryObservationPackageSetup(cfg, env)
	SnmpMetricPackageSetup(cfg, env)
	SignalPackageSetup(cfg, env)
}

func TsakCustomSetup(cfg *zygo.ZlispConfig, env *zygo.Zlisp) {
	log.Debug("DSL custom setup for TSAK is reached")
	callFun := `(def call _method)`
	_, err := env.EvalString(callFun)
	PanicOn(err)
	TsakGlobals(cfg, env)
}
