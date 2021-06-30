package dsl

import (
	"errors"

	. "github.com/glycerine/zygomys/zygo"
	"github.com/pieterclaerhout/go-log"
	"gonum.org/v1/gonum/interp"
)

type Predictor struct {
	Name string
	pc   interp.PiecewiseConstant
	pl   interp.PiecewiseLinear
	as   interp.AkimaSpline
	fb   interp.FritschButland
}

func (p *Predictor) SexpString(ps *PrintState) string {
	return p.Name
}

func (p *Predictor) Type() *RegisteredType {
	return GoStructRegistry.Registry[p.Name]
}

func (p *Predictor) Fit(x, y []float64) bool {
	err := errors.New("Generic predictions error")
	if len(x) != len(y) {
		log.Errorf("Fitting data must be of same arity and not %d x %d", len(x), len(y))
		return false
	}
	switch p.Name {
	case "PiecewiseConstant":
		err = p.pc.Fit(x, y)
	case "PiecewiseLinear":
		err = p.pl.Fit(x, y)
	case "AkimaSpline":
		err = p.as.Fit(x, y)
	case "FritschButland":
		err = p.fb.Fit(x, y)
	default:
		log.Errorf("We do not know how to use predictor %v", p.Name)
		return false
	}
	if err != nil {
		log.Errorf("Error fitting data for %s: %v", p.Name, err)
		return false
	}
	return true
}

func (p *Predictor) Predict(x float64) float64 {
	switch p.Name {
	case "PiecewiseConstant":
		return p.pc.Predict(x)
	case "PiecewiseLinear":
		return p.pl.Predict(x)
	case "AkimaSpline":
		return p.as.Predict(x)
	case "FritschButland":
		return p.fb.Predict(x)
	default:
		log.Errorf("We do not know how to use predictor %v", p.Name)
		return 0.0
	}
	return 0.0
}

func PredictorFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{}
}

func PredictorPackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def predictor (package "predictor"
     {
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}

func PredictorSetup() {
	GoStructRegistry.RegisterUserdef(&RegisteredType{GenDefMap: true, Factory: func(env *Zlisp, h *SexpHash) (interface{}, error) {
		return &Predictor{Name: "PiecewiseConstant"}, nil
	}}, true, "PiecewiseConstant")
	GoStructRegistry.RegisterUserdef(&RegisteredType{GenDefMap: true, Factory: func(env *Zlisp, h *SexpHash) (interface{}, error) {
		return &Predictor{Name: "PiecewiseLinear"}, nil
	}}, true, "PiecewiseLinear")
	GoStructRegistry.RegisterUserdef(&RegisteredType{GenDefMap: true, Factory: func(env *Zlisp, h *SexpHash) (interface{}, error) {
		return &Predictor{Name: "AkimaSpline"}, nil
	}}, true, "AkimaSpline")
	GoStructRegistry.RegisterUserdef(&RegisteredType{GenDefMap: true, Factory: func(env *Zlisp, h *SexpHash) (interface{}, error) {
		return &Predictor{Name: "FritschButland"}, nil
	}}, true, "FritschButland")
}
