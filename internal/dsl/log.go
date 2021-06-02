package dsl

import (
	"github.com/glycerine/zygomys/zygo"
	. "github.com/glycerine/zygomys/zygo"
	"github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/signal"
)

func LogFunction(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) < 1 {
		return SexpNull, WrongNargs
	}

	var str string

	switch expr := args[0].(type) {
	case *SexpStr:
		str = expr.S
	default:
		str = expr.SexpString(nil)
	}
	ar := make([]interface{}, len(args)-1)
	for i := 0; i < len(ar); i++ {
		switch x := args[i+1].(type) {
		case *SexpInt:
			ar[i] = x.Val
		case *SexpBool:
			ar[i] = x.Val
		case *SexpFloat:
			ar[i] = x.Val
		case *SexpChar:
			ar[i] = x.Val
		case *SexpStr:
			ar[i] = x.S
		case *SexpTime:
			ar[i] = x.Tm.In(NYC)
		default:
			ar[i] = args[i+1]
		}
	}
	switch name {
	case "log.Debug":
		if len(args) == 1 {
			log.Debug(str)
		} else {
			log.Debugf(str, ar...)
		}
	case "log.Warning":
		if len(args) == 1 {
			log.Warn(str)
		} else {
			log.Warnf(str, ar...)
		}
	case "log.Info":
		if len(args) == 1 {
			log.Info(str)
		} else {
			log.Infof(str, ar...)
		}
	case "log.Error":
		if len(args) == 1 {
			log.Error(str)
		} else {
			log.Errorf(str, ar...)
		}
	case "log.Fatal":
		signal.ExitRequest()
		if len(args) == 1 {
			log.Fatal(str)
		} else {
			log.Fatalf(str, ar...)
		}
	default:
		if len(args) == 1 {
			log.Info(str)
		} else {
			log.Infof(str, ar...)
		}
	}
	return SexpNull, nil
}

func LogFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"logdebug":   LogFunction,
		"logwarning": LogFunction,
		"loginfo":    LogFunction,
		"logerror":   LogFunction,
		"logfatal":   LogFunction,
	}
}

func TsakLogSetup(cfg *zygo.ZlispConfig, env *zygo.Zlisp) {
	logPkg := `(def log (package "log"
     { Debug := logdebug;
       Warning := logwarning;
       Info := loginfo;
       Error := logerror;
       Fatal := logfatal;
     }
  ))`
	_, err := env.EvalString(logPkg)
	PanicOn(err)
}
