package dsl

import (
	. "github.com/glycerine/zygomys/zygo"

	"github.com/vulogov/TSAK2/internal/conf"
)

func TsakGlobals(cfg *ZlispConfig, env *Zlisp) {
	env.AddGlobal("Name", &SexpStr{S: *conf.Name})
	env.AddGlobal("ID", &SexpStr{S: *conf.ID})
	env.AddGlobal("Mibsdb", &SexpStr{S: *conf.SNMPMibsdb})
	env.AddGlobal("Int32", &SexpInt{Val: 32})
	env.AddGlobal("Uint32", &SexpInt{Val: 33})
	env.AddGlobal("Int", &SexpInt{Val: 34})
	env.AddGlobal("Int64", &SexpInt{Val: 64})
	env.AddGlobal("Uint64", &SexpInt{Val: 65})
}
