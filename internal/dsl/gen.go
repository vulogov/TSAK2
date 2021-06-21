package dsl

import (
	. "github.com/glycerine/zygomys/zygo"
	"github.com/lrita/cmap"
	"github.com/pieterclaerhout/go-log"
	// "github.com/vulogov/TSAK2/internal/conf"
)

var DefaultHistorySize = 128

type GEN struct {
	Current     cmap.Cmap
	History     cmap.Cmap
	HistorySize int64
}

func GeneratorSetup() {
	log.Debug("Running Generator setup")
	GoStructRegistry.RegisterUserdef(&RegisteredType{GenDefMap: true, Factory: func(env *Zlisp, h *SexpHash) (interface{}, error) {
		return &GEN{HistorySize: int64(DefaultHistorySize)}, nil
	}}, true, "GEN")
}
