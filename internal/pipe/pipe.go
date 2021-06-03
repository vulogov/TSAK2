package pipe

import (
	. "github.com/glycerine/zygomys/zygo"
	"github.com/lrita/cmap"
	"github.com/pieterclaerhout/go-log"
)

var pipeSize = 1048576
var pipestore cmap.Cmap

func Createpipe(name string) bool {
	_, ok := pipestore.Load(name)
	if ok {
		return false
	}
	pipestore.Store(name, make(chan string, pipeSize))
	return true
}

func ZCreatePipe(env *Zlisp, name string,
	args []Sexp) (Sexp, error) {
	var pname string
	var n int64
	if len(args) >= 1 {
		for p := range args {
			switch expr := args[p].(type) {
			case *SexpStr:
				pname = expr.S
			default:
				pname = expr.SexpString(nil)
			}
			log.Debugf("Creating chan: %v", pname)
			if Createpipe(pname) {
				n += 1
			}
		}
	} else {
		return SexpNull, WrongNargs
	}
	return &SexpInt{Val: n}, nil
}

func PipeLen(name string) int64 {
	c, ok := pipestore.Load(name)
	if ok {
		return int64(len(c.(chan string)))
	}
	return 0
}

func ZPipeLen(env *Zlisp, name string,
	args []Sexp) (Sexp, error) {
	var pname string
	if len(args) == 1 {
		switch expr := args[0].(type) {
		case *SexpStr:
			pname = expr.S
		default:
			pname = expr.SexpString(nil)
		}
		return &SexpInt{Val: PipeLen(pname)}, nil
	} else {
		return SexpNull, WrongNargs
	}
}

func Sendpipe(name string, data string) int {
	c, ok := pipestore.Load(name)
	if ok {
		c.(chan string) <- data
		return len(c.(chan string))
	}
	return 0
}

func ZSendPipe(env *Zlisp, name string,
	args []Sexp) (Sexp, error) {
	var pname string
	var data string
	if len(args) >= 1 {
		switch expr := args[0].(type) {
		case *SexpStr:
			pname = expr.S
		default:
			pname = expr.SexpString(nil)
		}
		for p := range args[1:] {
			switch expr := args[p+1].(type) {
			case *SexpStr:
				data = expr.S
			default:
				data = expr.SexpString(nil)
			}
			log.Debugf("Sending %v bytes to channel %v", len(data), pname)
			Sendpipe(pname, data)
		}
	} else {
		return SexpNull, WrongNargs
	}
	return &SexpInt{Val: PipeLen(pname)}, nil
}

func Recvpipe(name string) (res string) {
	c, ok := pipestore.Load(name)
	if ok {
		if len(c.(chan string)) > 0 {
			res = <-c.(chan string)
			return
		}
	}
	return
}

func ZRecvpipe(env *Zlisp, name string,
	args []Sexp) (Sexp, error) {
	var pname string
	if len(args) == 1 {
		switch expr := args[0].(type) {
		case *SexpStr:
			pname = expr.S
		default:
			pname = expr.SexpString(nil)
		}
		return &SexpStr{S: Recvpipe(pname)}, nil
	} else {
		return SexpNull, WrongNargs
	}
}

func PipeFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"pipecreate": ZCreatePipe,
		"pipesend":   ZSendPipe,
		"pipelen":    ZPipeLen,
		"piperecv":   ZRecvpipe,
	}
}
