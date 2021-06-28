package dsl

import (
	"fmt"

	. "github.com/glycerine/zygomys/zygo"
	floats "gonum.org/v1/gonum/floats"
)

func FloatLogSpan(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) != 3 {
		return SexpNull, WrongNargs
	}
	if !IsInt(args[0]) {
		return SexpNull, fmt.Errorf("Arity of array in %s must be integer", name)
	}
	if !IsFloat(args[1]) {
		return SexpNull, fmt.Errorf("Lower in %s must be float", name)
	}
	if !IsFloat(args[2]) {
		return SexpNull, fmt.Errorf("Upper in %s must be float", name)
	}
	res := make([]float64, int(AsAny(args[0]).(int64)))
	floats.LogSpan(res, float64(AsAny(args[1]).(float64)), float64(AsAny(args[2]).(float64)))
	return ArrayofFloatsToLispArray(env, res), nil
}

func FloatMath(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	var fres []float64
	if len(args) < 2 {
		return SexpNull, WrongNargs
	}
	if IsArray(args[0]) {
		fres = ArrayofFloatsToArray(args[0])
		fmt.Println(fres)
	} else {
		return SexpNull, fmt.Errorf("Invalid data in %s operation", name)
	}
	for _, v := range args[1:] {
		switch e := v.(type) {
		case *SexpArray:
			if len(e.Val) != len(fres) {
				return SexpNull, fmt.Errorf("Arity of arrays in %s must match", name)
			}
			a := ArrayofFloatsToArray(e)
			switch name {
			case "float.Add":
				floats.Add(fres, a)
			case "float.Div":
				floats.Div(fres, a)
			case "float.Mul":
				floats.Mul(fres, a)
			case "float.Dot":
				floats.Dot(fres, a)
			}
		}
	}
	return ArrayofFloatsToLispArray(env, fres), nil
}

func FloatTo(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	res := &SexpArray{Val: make([]Sexp, 0), Env: env}
	if len(args) == 0 {
		return SexpNull, WrongNargs
	}
	for _, v := range args {
		switch e := v.(type) {
		case *SexpInt:
			res.Val = append(res.Val, &SexpFloat{Val: float64(e.Val)})
		case *SexpFloat:
			res.Val = append(res.Val, &SexpFloat{Val: float64(e.Val)})
		case *SexpArray:
			for _, v1 := range e.Val {
				switch e1 := v1.(type) {
				case *SexpInt:
					res.Val = append(res.Val, &SexpFloat{Val: float64(e1.Val)})
				case *SexpFloat:
					res.Val = append(res.Val, &SexpFloat{Val: float64(e1.Val)})
				default:
					return SexpNull, fmt.Errorf("Invalid data in float.To inner array conversion")
				}
			}
		default:
			return SexpNull, fmt.Errorf("Invalid data in float.To conversion")
		}
	}
	return res, nil
}

func FloatBytes(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	res := make([]float64, 0)
	if len(args) == 0 {
		return SexpNull, WrongNargs
	}
	for _, v := range args {
		switch e := v.(type) {
		case *SexpInt:
			res = append(res, float64(e.Val))
		case *SexpFloat:
			res = append(res, float64(e.Val))
		case *SexpArray:
			for _, v1 := range e.Val {
				switch e1 := v1.(type) {
				case *SexpInt:
					res = append(res, float64(e1.Val))
				case *SexpFloat:
					res = append(res, float64(e1.Val))
				default:
					return SexpNull, fmt.Errorf("Invalid data in float.Bytes inner array conversion")
				}
			}
		default:
			return SexpNull, fmt.Errorf("Invalid data in float.Bytes conversion")
		}
	}
	switch name {
	case "float.Bytes":
		return ArrayofFloatsToLispArray(env, res), nil
	case "float.KBytes":
		return ArrayofFloatsToLispArray(env, ArrayOfFloatsMulOn(res, 1024)), nil
	case "float.MBytes":
		return ArrayofFloatsToLispArray(env, ArrayOfFloatsMulOn(res, (1024*1024))), nil
	case "float.GBytes":
		return ArrayofFloatsToLispArray(env, ArrayOfFloatsMulOn(res, (1024*1024*1024))), nil
	case "float.Int":
		return ArrayofFloatsToIntLispArray(env, res), nil
	}
	return SexpNull, fmt.Errorf("I do not know how to perform this conversion: %s", name)
}

func FloatFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"floatto":      FloatTo,
		"floatadd":     FloatMath,
		"floatdiv":     FloatMath,
		"floatmul":     FloatMath,
		"floatdot":     FloatMath,
		"floatlogspan": FloatLogSpan,
		"floatbytes":   FloatBytes,
		"floatkbytes":  FloatBytes,
		"floatmbytes":  FloatBytes,
		"floatgbytes":  FloatBytes,
		"floatint":     FloatBytes,
	}
}

func FloatPackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def float (package "float"
     { To := floatto;
			 Add := floatadd;
			 Div := floatdiv;
			 Mul := floatmul;
			 Dot := floatdot;
			 LogSpan := floatlogspan;
			 Bytes := floatbytes;
			 KBytes := floatkbytes;
			 MBytes := floatmbytes;
			 GBytes := floatgbytes;
			 Int := floatint;
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}
