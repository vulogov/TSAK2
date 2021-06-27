package dsl

import (
	. "github.com/glycerine/zygomys/zygo"
)

func RandFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{}
}

func RandPackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def rand (package "rand"
     {
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}
