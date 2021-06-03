package dsl

import (
	"fmt"

	. "github.com/glycerine/zygomys/zygo"
	"github.com/goml/gobrain"
	"github.com/pieterclaerhout/go-log"
)

type Perceptron struct {
	f            gobrain.FeedForward
	Name         string `json:"Name" msg:"Name"`
	isConfigured bool
	isTrained    bool
	Epoch        int `json:"Epoch" msg:"Epoch"`
	isReport     bool
	Lrate        float64 `json:"Lrate" msg:"Lrate"`
	Momentum     float64 `json:"Momentum" msg:"Momentum"`
	In           int     `json:"In" msg:"In"`
	Hidden       int     `json:"Hidden" msg:"Hidden"`
	Out          int     `json:"Out" msg:"Out"`
}

func (ff *Perceptron) Setup(n string, input, hidden, output int) (res bool, err error) {
	log.Debugf("Configuring Perceptron[%v]", n)
	ff.f = gobrain.FeedForward{}
	ff.isConfigured = false
	ff.isTrained = false
	ff.isReport = false
	ff.In = input
	ff.Hidden = hidden
	ff.Out = output
	ff.Name = n
	ff.f.Init(input, hidden, output)
	ff.Epoch = 1000
	ff.Lrate = 0.6
	ff.Momentum = 0.4
	ff.isReport = false
	ff.isConfigured = true
	res = true
	return
}

func (ff *Perceptron) Configure(epoch int, lrate float64, momentum float64, isReport bool) {
	ff.Epoch = epoch
	ff.Lrate = lrate
	ff.Momentum = momentum
	ff.isReport = isReport
	ff.isConfigured = true
}

func (ff *Perceptron) DisplayPerceptron(from string) {
	fmt.Printf("Perceptron: %s %#v", from, ff)
}

func PerceptronModuleFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{}
}

func PerceptronSetup() {
	log.Debug("Running perceptron setup")
	GoStructRegistry.RegisterUserdef(&RegisteredType{GenDefMap: true, Factory: func(env *Zlisp, h *SexpHash) (interface{}, error) {
		return &Perceptron{}, nil
	}}, true, "Perceptron")
}
