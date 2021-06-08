package dsl

import (
	. "github.com/glycerine/zygomys/zygo"
)

func AsString(expr Sexp) string {
	switch e := expr.(type) {
	case *SexpStr:
		return e.S
	default:
		return expr.SexpString(nil)
	}
}
