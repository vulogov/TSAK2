package dsl

import (
	"github.com/glycerine/zygomys/zygo"
	. "github.com/glycerine/zygomys/zygo"
	"github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/pipe"
)

var Env *Zlisp
var Cfg *ZlispConfig

func GetTheAnswer(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	log.Debug("Someone is looking for an answer. Well it is 42")
	return &SexpInt{Val: int64(42)}, nil
}

func TsakBuiltinFunctions() map[string]zygo.ZlispUserFunction {
	log.Debug("Registering TSAK built-in functions")
	return MergeFuncMap(
		AllBuiltinFunctions(),
		AllTsakCoreFunctions(),
		LogFunctions(),
		PerceptronModuleFunctions(),
		TelemetryObservationFunctions(),
		SnmpMetricFunctions(),
		MIBSFunctions(),
		pipe.PipeFunctions(),
		SignalFunctions(),
		GeneratorFunctions(),
		SleepFunctions(),
		NowFunctions(),
		FakeFunctions(),
	)
}

func AllTsakCoreFunctions() map[string]ZlispUserFunction {
	log.Debug("Registering TSAK core functions")
	return map[string]ZlispUserFunction{
		"answer": GetTheAnswer,
	}
}

func AllEnvInitBeforeCreationOfEnv() {
	log.Debug("DSL initialization before environment creation")
	PerceptronSetup()
	TelemetryObservationMatrixSetup()
	MIBSDatatypeSetup()
	GeneratorSetup()
}
