package dsl

import (
	. "github.com/glycerine/zygomys/zygo"
	"github.com/pieterclaerhout/go-log"
)

func TsakPipeSetup(cfg *ZlispConfig, env *Zlisp) {
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

func TsakStandardSetup(cfg *ZlispConfig, env *Zlisp) {
	log.Debug("DSL setup for TSAK is reached")
	TsakLogSetup(cfg, env)
	TsakPipeSetup(cfg, env)
	TelemetryObservationPackageSetup(cfg, env)
	SnmpMetricPackageSetup(cfg, env)
	SignalPackageSetup(cfg, env)
	GeneratorPackageSetup(cfg, env)
	SleepPackageSetup(cfg, env)
	NowPackageSetup(cfg, env)
	FakePackageSetup(cfg, env)
	MathPackageSetup(cfg, env)
	RandPackageSetup(cfg, env)
	NormalizePackageSetup(cfg, env)
	FloatPackageSetup(cfg, env)
	PredictorPackageSetup(cfg, env)
}

func TsakCustomSetup(cfg *ZlispConfig, env *Zlisp) {
	log.Debug("DSL custom setup for TSAK is reached")
	callFun := `(def call _method)`
	_, err := env.EvalString(callFun)
	PanicOn(err)
	TsakGlobals(cfg, env)
	GenerateMetricStart()
}
