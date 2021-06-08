package dsl

import (
	. "github.com/glycerine/zygomys/zygo"
)

func ArrayofStringsToArray(expr Sexp) (res []string) {
	switch e := expr.(type) {
	case *SexpArray:
		res = make([]string, len(e.Val))
		for n, v := range e.Val {
			switch z := v.(type) {
			case *SexpStr:
				res[n] = z.S
			default:
				res[n] = v.SexpString(nil)
			}
		}
		return
	}
	res = make([]string, 0)
	return
}
