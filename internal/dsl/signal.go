package dsl

import (
	"fmt"

	. "github.com/glycerine/zygomys/zygo"
	"github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/signal"
)

func SignalExitRequest(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	signal.ExitRequest()
	return &SexpInt{Val: int64(signal.Len())}, nil
}
func SignalExitRequested(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	res := signal.ExitRequested()
	return &SexpBool{Val: res}, nil
}
func SignalReserve(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	var n int64
	n = 0
	if len(args) == 0 {
		n = 1
	} else {
		if !IsInt(args[0]) {
			return SexpNull, fmt.Errorf("First argument must be string")
		}
		switch v := args[0].(type) {
		case *SexpInt:
			n = v.Val
		default:
			n = 1
		}
	}
	log.Debugf("Reserving %d goroutines", n)
	res := signal.Reserve(int(n))
	return &SexpInt{Val: int64(res)}, nil
}
func SignalRelease(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	var n int64
	n = 0
	if len(args) == 0 {
		n = 1
	} else {
		if !IsInt(args[0]) {
			return SexpNull, fmt.Errorf("First argument must be string")
		}
		switch v := args[0].(type) {
		case *SexpInt:
			n = v.Val
		default:
			n = 1
		}
	}
	log.Debugf("Reserving %d goroutines", n)
	res := signal.Release(int(n))
	return &SexpInt{Val: int64(res)}, nil
}

func SignalFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"signalexitrequest":   SignalExitRequest,
		"signalexitrequested": SignalExitRequested,
		"signalreserve":       SignalReserve,
		"signalrelease":       SignalRelease,
		"quit":                SignalExitRequest,
	}
}

func SignalPackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def signal (package "signal"
     { ExitRequest := signalexitrequest;
       ExitRequested := signalexitrequested;
       Reserve := signalreserve;
			 Release := signalrelease;
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}
