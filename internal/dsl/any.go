package dsl

import (
	. "github.com/glycerine/zygomys/zygo"
)

func AsAny(expr Sexp) interface{} {
	switch e := expr.(type) {
	case *SexpFloat:
		return e.Val
	case *SexpInt:
		return e.Val
	case *SexpStr:
		return e.S
	}
	return expr.SexpString(nil)
}
