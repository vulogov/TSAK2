package dsl

import (
	"fmt"
	"math"

	. "github.com/glycerine/zygomys/zygo"
)

func MathTrig(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	var x float64
	if len(args) == 0 {
		return SexpNull, WrongNargs
	}
	switch e := args[0].(type) {
	case *SexpInt:
		x = float64(e.Val)
	case *SexpFloat:
		x = float64(e.Val)
	default:
		return SexpNull, fmt.Errorf("First argument must be int or float")
	}
	switch name {
	case "math.Sin":
		return &SexpFloat{Val: math.Sin(x)}, nil
	case "math.Cos":
		return &SexpFloat{Val: math.Cos(x)}, nil
	case "math.Abs":
		return &SexpFloat{Val: math.Abs(x)}, nil
	case "math.Log10":
		return &SexpFloat{Val: math.Log10(x)}, nil
	case "math.Log2":
		return &SexpFloat{Val: math.Log2(x)}, nil
	case "math.Exp":
		return &SexpFloat{Val: math.Exp(x)}, nil
	case "math.Pow10":
		return &SexpFloat{Val: math.Pow10(int(x))}, nil
	}
	return SexpNull, fmt.Errorf("Requested math computation can not be performed: %v", name)
}

func MathFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"mathsin":   MathTrig,
		"mathcos":   MathTrig,
		"mathabs":   MathTrig,
		"mathpow10": MathTrig,
		"mathlog10": MathTrig,
		"mathlog2":  MathTrig,
		"mathexp":   MathTrig,
	}
}

func MathPackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def math (package "math"
     { Sin := mathsin;
       Cos := mathcos;
       Abs := mathabs;
       Log10 := mathlog10;
       Log2 := mathlog2;
       Exp := mathexp;
       Pow10 := mathpow10;
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}
