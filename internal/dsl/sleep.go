package dsl

import (
	"time"

	. "github.com/glycerine/zygomys/zygo"
)

func SleepSecond(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	var n time.Duration
	if len(args) == 0 {
		n = 1
	} else {
		switch e := args[0].(type) {
		case *SexpInt:
			n = time.Duration(e.Val)
		default:
			n = 1
		}
	}
	time.Sleep(n * time.Second)
	return &SexpInt{Val: int64(time.Now().UTC().UnixNano())}, nil
}

func SleepMillisecond(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	var n time.Duration
	if len(args) == 0 {
		n = 100
	} else {
		switch e := args[0].(type) {
		case *SexpInt:
			n = time.Duration(e.Val)
		default:
			n = 100
		}
	}
	time.Sleep(n * time.Millisecond)
	return &SexpInt{Val: int64(time.Now().UTC().UnixNano())}, nil
}

func SleepFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"sleepsecond":      SleepSecond,
		"sleepmillisecond": SleepMillisecond,
	}
}

func SleepPackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def sleep (package "sleep"
     { Second := sleepsecond;
       Millisecond := sleepmillisecond;
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}
