package dsl

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/glycerine/zygomys/zygo"
	distuv "gonum.org/v1/gonum/stat/distuv"
)

var RandSource = rand.NewSource(int64(time.Now().UTC().UnixNano()))
var Rand = rand.New(RandSource)

func randFloatsToFloat64(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = float64(min + rand.Float64()*(max-min))
	}
	return res
}

func RandRand(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) == 0 {
		switch name {
		case "rand.FloatMax":
			return &SexpFloat{Val: float64(Rand.ExpFloat64())}, nil
		case "rand.Float":
			return &SexpFloat{Val: float64(Rand.Float64())}, nil
		case "rand.FloatNorm":
			return &SexpFloat{Val: float64(Rand.NormFloat64())}, nil
		case "rand.Int":
			return &SexpInt{Val: int64(Rand.Int63())}, nil
		default:
			return SexpNull, fmt.Errorf("Do not know how to compute: %v", name)
		}
	}
	if len(args) == 2 && (name == "rand.Intn" || name == "rand.Fltn") {
		min := int(0)
		switch e := args[0].(type) {
		case *SexpInt:
			min = int(e.Val)
		default:
			return SexpNull, fmt.Errorf("First parameter must be integer")
		}
		max := int(0)
		switch e1 := args[1].(type) {
		case *SexpInt:
			max = int(e1.Val)
		default:
			return SexpNull, fmt.Errorf("Second parameter must be integer")
		}
		v := (max - min) + 1
		if name != "rand.Fltn" {
			return &SexpInt{Val: int64(rand.Intn(int(v)) + min)}, nil
		} else {
			return &SexpFloat{Val: float64(rand.Intn(int(v)) + min)}, nil
		}
	}
	return SexpNull, fmt.Errorf("Do not know how to compute: %v", name)
}

func RandToFloat(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) == 2 {
		min := float64(0.0)
		switch e := args[0].(type) {
		case *SexpFloat:
			min = float64(e.Val)
		default:
			return SexpNull, fmt.Errorf("First parameter must be float")
		}
		max := float64(0.0)
		switch e1 := args[1].(type) {
		case *SexpFloat:
			max = float64(e1.Val)
		default:
			return SexpNull, fmt.Errorf("Second parameter must be float")
		}
		val := randFloatsToFloat64(min, max, 1)
		return &SexpFloat{Val: val[0]}, nil
	}
	return SexpNull, fmt.Errorf("Do not know how to compute: %v", name)
}

func RandFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"randfloatmax":  RandRand,
		"randfloat":     RandRand,
		"randfloatnorm": RandRand,
		"randfloatint":  RandRand,
		"randintintn":   RandRand,
		"randintintnfl": RandRand,
		"randnfloat":    RandToFloat,
	}
}

func RandPackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def rand (package "rand"
     { FloatMax := randfloatmax;
			 Float := randfloat;
			 FloatNorm := randfloatnorm;
			 Int := randfloatint;
			 Intn := randintintn;
			 Fltn := randintintnfl ;
			 NFloat := randnfloat ;
     }
  ))`
	_, err := env.EvalString(myPkg)
	rand.Seed(time.Now().UnixNano())
	PanicOn(err)
}

func RandomSetup() {
	GoStructRegistry.RegisterUserdef(&RegisteredType{GenDefMap: true, Factory: func(env *Zlisp, h *SexpHash) (interface{}, error) {
		return &distuv.Beta{}, nil
	}}, true, "BetaRandom")
}
