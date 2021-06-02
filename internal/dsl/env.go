package dsl

import (
	"github.com/glycerine/zygomys/zygo"
	"github.com/pieterclaerhout/go-log"
)

func TsakStandardSetup(cfg *zygo.ZlispConfig, env *zygo.Zlisp) {
	log.Debug("DSL setup for TSAK is reached")
	TsakLogSetup(cfg, env)
}

func TsakCustomSetup(cfg *zygo.ZlispConfig, env *zygo.Zlisp) {
	log.Debug("DSL custom setup for TSAK is reached")
	callFun := `(def call _method)`
	_, err := env.EvalString(callFun)
	PanicOn(err)
}
