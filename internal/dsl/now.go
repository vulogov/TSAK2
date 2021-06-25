package dsl

import (
	"time"

	. "github.com/glycerine/zygomys/zygo"
)

func NowUTC(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	return &SexpInt{Val: int64(time.Now().UTC().Unix())}, nil
}

func NowUTCNano(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	return &SexpInt{Val: int64(time.Now().UTC().UnixNano())}, nil
}

func NowLocal(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	return &SexpInt{Val: int64(time.Now().Unix())}, nil
}

func NowLocalNano(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	return &SexpInt{Val: int64(time.Now().UnixNano())}, nil
}

func NowFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"nowutc":       NowUTC,
		"nowutcnano":   NowUTCNano,
		"nowlocal":     NowLocal,
		"nowlocalnano": NowLocalNano,
	}
}

func NowPackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def now (package "now"
     { UTC := nowutc;
       UTCNano := nowutcnano;
			 Local := nowlocal;
	     LocalNano := nowlocalnano;
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}
