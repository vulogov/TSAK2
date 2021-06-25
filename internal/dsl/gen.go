package dsl

import (
	"fmt"
	re "regexp"
	"sync"
	"time"

	"github.com/edwingeng/deque"
	. "github.com/glycerine/zygomys/zygo"
	"github.com/lrita/cmap"
	"github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/conf"
	"github.com/vulogov/TSAK2/internal/signal"
)

var DefaultHistorySize = int64(128)
var Gen = GEN{HistorySize: DefaultHistorySize}
var GenInterval int64
var PipelineInterval int64

type GEN struct {
	Current     cmap.Cmap
	History     cmap.Cmap
	X           cmap.Cmap
	Gen         cmap.Cmap
	Default     interface{}
	HistorySize int64
	Lock        sync.Mutex
}

func (gen *GEN) _create(source, name string, generator string, defval interface{}) bool {
	log.Debugf("Creating GEN for %s@%s %v", source, name, gen.HistorySize)
	key := makekey(source, name)
	_, ok := gen.Current.Load(key)
	if ok {
		return false
	}
	gen.Default = defval
	gen.Current.Store(key, defval)
	gen.Gen.Store(key, generator)
	gen.History.Store(key, deque.NewDeque())
	gen.X.Store(key, deque.NewDeque())
	return true
}

func (gen *GEN) MakeInt(source, name string, generator string, val int64) bool {
	return gen._create(source, name, generator, val)
}

func (gen *GEN) MakeFloat(source, name string, generator string, val float64) bool {
	return gen._create(source, name, generator, val)
}

func (gen *GEN) MakeString(source, name string, generator string, val string) bool {
	return gen._create(source, name, generator, val)
}

func (gen *GEN) Get(source, name string) interface{} {
	key := makekey(source, name)
	h, ok := gen.History.Load(key)
	if !ok {
		log.Errorf("GEN: Unable to locate history for %v", key)
		return nil
	}
	if h.(deque.Deque).Len() == 0 {
		gen.Current.Store(key, gen.Default)
	}
	d, ok := gen.Current.Load(key)
	if !ok {
		log.Errorf("GEN: Unable to locate current value for %v", key)
		return nil
	}
	return d
}

func (gen *GEN) Take(source, name string) interface{} {
	key := makekey(source, name)
	h, ok := gen.History.Load(key)
	if !ok {
		log.Errorf("GEN: Unable to locate history for %v", key)
		return nil
	}
	x, ok := gen.X.Load(key)
	if !ok {
		log.Errorf("GEN: Unable to locate X for %v", key)
		return nil
	}
	if h.(deque.Deque).Len() == 0 {
		gen.Current.Store(key, gen.Default)
	} else {
		gen.Current.Store(key, h.(deque.Deque).PopFront())
		x.(deque.Deque).PopFront()
	}
	d, ok := gen.Current.Load(key)
	if !ok {
		log.Errorf("GEN: Unable to locate current value for %v", key)
		return nil
	}
	return d
}

func (gen *GEN) GetInt(source, name string) int64 {
	res := gen.Get(source, name)
	switch e := res.(type) {
	case int64:
		return int64(e)
	}
	return int64(0)
}

func (gen *GEN) GetFloat(source, name string) float64 {
	return float64(gen.Get(source, name).(float64))
}

func (gen *GEN) GetString(source, name string) string {
	return string(gen.Get(source, name).(string))
}

func (gen *GEN) trim(source, name string) int {
	res := 0
	key := makekey(source, name)
	h, ok := gen.History.Load(key)
	if ok {
		for h.(deque.Deque).Len() > int(gen.HistorySize) {
			h.(deque.Deque).PopFront()
			res += 1
		}
	}
	x, ok := gen.X.Load(key)
	if ok {
		for x.(deque.Deque).Len() > int(gen.HistorySize) {
			x.(deque.Deque).PopFront()
		}
	}
	return res
}

func (gen *GEN) Compute(source, name string) interface{} {
	var res interface{}
	key := makekey(source, name)
	g, ok := gen.Gen.Load(key)
	if !ok {
		return nil
	}
	log.Debugf("Applying %v to %v", g, key)
	env := Env.Clone()
	env.AddGlobal("Source", &SexpStr{S: source})
	env.AddGlobal("Key", &SexpStr{S: name})
	r, err := env.EvalString(g.(string))
	switch err {
	case nil:
	case NoExpressionsFound:
		log.Errorf("Error: %v", err)
		Env.Clear()
		return nil
	default:
		log.Errorf("Error: %v", err)
		fmt.Print(Env.GetStackTrace(err))
		Env.Clear()
		return nil
	}
	switch e := r.(type) {
	case *SexpArray:
		res = AsAny(e.Val[len(e.Val)-1])
	default:
		res = AsAny(e)
	}
	gen.Current.Store(key, res)
	h, ok := gen.History.Load(key)
	if !ok {
		log.Debugf("Can not find history for %s@%s", source, name)
		return nil
	}
	h.(deque.Deque).PushBack(res)
	log.Debugf("History size for %v is %d", key, h.(deque.Deque).Len())
	x, ok := gen.X.Load(key)
	if !ok {
		log.Debugf("Can not find X for %s@%s", source, name)
		return nil
	}
	x.(deque.Deque).PushBack(time.Now().UTC().UnixNano())
	gen.trim(source, name)
	env.Clear()
	return res
}

func (gen *GEN) Len(source, name string) int64 {
	key := makekey(source, name)
	h, ok := gen.History.Load(key)
	if !ok {
		return int64(0)
	}
	return int64(h.(deque.Deque).Len())
}

func (gen *GEN) GetHistory(source, name string) []interface{} {
	res := make([]interface{}, 0)
	key := makekey(source, name)
	h, ok := gen.History.Load(key)
	if ok {
		h.(deque.Deque).Range(func(i int, v interface{}) bool {
			res = append(res, v)
			return true
		})
	}
	return res
}

func (gen *GEN) GetX(source, name string) []interface{} {
	res := make([]interface{}, 0)
	key := makekey(source, name)
	x, ok := gen.X.Load(key)
	if ok {
		x.(deque.Deque).Range(func(i int, v interface{}) bool {
			res = append(res, v)
			return true
		})
	}
	return res
}

func GeneratorMake(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) != 4 {
		return SexpNull, WrongNargs
	}
	if !IsArray(args[0]) {
		return SexpNull, fmt.Errorf("First argument must be array")
	}
	if !IsArray(args[1]) {
		return SexpNull, fmt.Errorf("Second argument must be array")
	}
	if !IsString(args[2]) {
		return SexpNull, fmt.Errorf("Second argument must be array")
	}
	ar1 := ArrayofStringsToArray(args[0])
	ar2 := ArrayofStringsToArray(args[1])

	n := 0
	for _, s := range ar1 {
		for _, k := range ar2 {
			switch v := args[3].(type) {
			case *SexpFloat:
				Gen.MakeFloat(s, k, AsString(args[2]), v.Val)
			case *SexpInt:
				Gen.MakeInt(s, k, AsString(args[2]), v.Val)
			case *SexpStr:
				Gen.MakeString(s, k, AsString(args[2]), v.S)
			default:
				Gen.MakeString(s, k, AsString(args[2]), args[3].SexpString(nil))
			}
			n += 1
		}
	}
	return &SexpInt{Val: int64(n)}, nil
}

func GeneratorCompute(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) != 2 {
		return SexpNull, WrongNargs
	}
	if !IsArray(args[0]) {
		return SexpNull, fmt.Errorf("First argument must be array")
	}
	if !IsArray(args[1]) {
		return SexpNull, fmt.Errorf("Second argument must be array")
	}

	ar1 := ArrayofStringsToArray(args[0])
	ar2 := ArrayofStringsToArray(args[1])

	n := 0
	for _, s := range ar1 {
		for _, k := range ar2 {
			Gen.Compute(s, k)
			n += 1
		}
	}
	return &SexpInt{Val: int64(n)}, nil
}

func GeneratorHistory(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	var res = &SexpArray{Env: env}
	if len(args) != 2 {
		return SexpNull, WrongNargs
	}
	if !IsString(args[0]) {
		return SexpNull, fmt.Errorf("First argument must be string")
	}
	if !IsString(args[1]) {
		return SexpNull, fmt.Errorf("Second argument must be string")
	}
	h := Gen.GetHistory(AsString(args[0]), AsString(args[1]))
	for _, v := range h {
		switch e := v.(type) {
		case int64, int:
			res.Val = append(res.Val, &SexpInt{Val: int64(e.(int64))})
		case float64:
			res.Val = append(res.Val, &SexpFloat{Val: float64(e)})
		case string:
			res.Val = append(res.Val, &SexpStr{S: string(e)})
		default:
			continue
		}
	}
	return res, nil
}

func GeneratorX(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	var res = &SexpArray{Env: env}
	if len(args) != 2 {
		return SexpNull, WrongNargs
	}
	if !IsString(args[0]) {
		return SexpNull, fmt.Errorf("First argument must be string")
	}
	if !IsString(args[1]) {
		return SexpNull, fmt.Errorf("Second argument must be string")
	}
	h := Gen.GetX(AsString(args[0]), AsString(args[1]))
	for _, v := range h {
		switch e := v.(type) {
		case int64, int:
			res.Val = append(res.Val, &SexpInt{Val: int64(e.(int64))})
		default:
			continue
		}
	}
	return res, nil
}

func GeneratorGet(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) != 2 {
		return SexpNull, WrongNargs
	}
	if !IsString(args[0]) {
		return SexpNull, fmt.Errorf("First argument must be string")
	}
	if !IsString(args[1]) {
		return SexpNull, fmt.Errorf("Second argument must be string")
	}
	h := Gen.Get(AsString(args[0]), AsString(args[1]))
	if h != nil {
		switch e := h.(type) {
		case int64, int:
			return &SexpInt{Val: int64(e.(int64))}, nil
		case float64:
			return &SexpFloat{Val: float64(e)}, nil
		case string:
			return &SexpStr{S: string(e)}, nil
		}
	}
	return SexpNull, fmt.Errorf("Gen.Get returns nil for %v %v", AsString(args[0]), AsString(args[1]))
}

func GeneratorTake(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) != 2 {
		return SexpNull, WrongNargs
	}
	if !IsString(args[0]) {
		return SexpNull, fmt.Errorf("First argument must be string")
	}
	if !IsString(args[1]) {
		return SexpNull, fmt.Errorf("Second argument must be string")
	}
	h := Gen.Take(AsString(args[0]), AsString(args[1]))
	if h != nil {
		switch e := h.(type) {
		case int64, int:
			return &SexpInt{Val: int64(e.(int64))}, nil
		case float64:
			return &SexpFloat{Val: float64(e)}, nil
		case string:
			return &SexpStr{S: string(e)}, nil
		}
	}
	return SexpNull, fmt.Errorf("Gen.Take returns nil for %v %v", AsString(args[0]), AsString(args[1]))
}

func GeneratorSetInterval(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) == 0 {
		GenInterval = int64(*conf.GenInterval)
	} else {
		switch e := args[0].(type) {
		case *SexpInt:
			GenInterval = int64(e.Val)
		}
	}
	log.Debugf("Telemetry greneration loop will be engaged every %v seconds", GenInterval)
	return &SexpInt{Val: int64(GenInterval)}, nil
}

func GeneratorPipelineInterval(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) == 0 {
		PipelineInterval = int64(*conf.GenInterval)
	} else {
		switch e := args[0].(type) {
		case *SexpInt:
			PipelineInterval = int64(e.Val)
		}
	}
	log.Debugf("Telemetry greneration loop will be engaged every %v seconds", PipelineInterval)
	return &SexpInt{Val: int64(PipelineInterval)}, nil
}

func GenerateMetricLoop() {
	log.Debug("Telemetry generation loop is started")
	GenInterval = int64(*conf.GenInterval)
	log.Debugf("Telemetry generation loop will be engaged every %v seconds", GenInterval)
	rekey := re.MustCompile(`\@`)
	for !signal.ExitRequested() {
		time.Sleep(time.Duration(GenInterval) * time.Second)
		log.Debug("Generator computing loop started")
		Gen.Gen.Range(func(key, val interface{}) bool {
			sk := rekey.Split(string(key.(string)), 2)
			if len(sk) != 2 {
				return true
			}
			log.Debug("Telemetry Generator LOCK engaged")
			Gen.Lock.Lock()
			Gen.Compute(sk[0], sk[1])
			log.Debug("Telemetry Generator UNLOCK engaged")
			Gen.Lock.Unlock()
			return true
		})
	}
	log.Debug("Telemetry generation loop is terminated")
}

func PipelineMetricLoop() {
	log.Debug("Telemetry pipeline loop is started")
	PipelineInterval = int64(*conf.GenInterval)
	log.Debugf("Telemetry pipeline loop will be engaged every %v seconds", GenInterval)
	rekey := re.MustCompile(`\@`)
	for !signal.ExitRequested() {
		time.Sleep(time.Duration(PipelineInterval) * time.Second)
		log.Debug("Pipeline loop started")
		log.Debug("Telemetry Pipeline LOCK engaged")
		Gen.Lock.Lock()
		Gen.Gen.Range(func(key, val interface{}) bool {
			sk := rekey.Split(string(key.(string)), 2)
			if len(sk) != 2 {
				return true
			}
			res := Gen.Get(sk[0], sk[1])
			switch e := res.(type) {
			case int, int64:
				TOM.AddInt(sk[0], sk[1], int64(e.(int64)))
			case float64:
				TOM.AddFloat(sk[0], sk[1], float64(e))
			case string:
				TOM.AddString(sk[0], sk[1], string(e))
			}
			return true
		})
		log.Debug("Telemetry Pipeline UNLOCK engaged")
		Gen.Lock.Unlock()
	}
	log.Debug("Telemetry pipeline loop is terminated")
}

func GenerateMetricStart() {
	log.Debug("Starting telemetry generation and delivery")
	go GenerateMetricLoop()
	go PipelineMetricLoop()
}

func GeneratorFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"generatormake":             GeneratorMake,
		"generatorcompute":          GeneratorCompute,
		"generatorhistory":          GeneratorHistory,
		"generatorx":                GeneratorX,
		"generatorget":              GeneratorGet,
		"generatortake":             GeneratorTake,
		"generatorsetinterval":      GeneratorSetInterval,
		"generatorpipelineinterval": GeneratorPipelineInterval,
	}
}

func GeneratorPackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def generator (package "generator"
     { Make := generatormake;
       Compute := generatorcompute;
       History := generatorhistory;
       X := generatorx;
       Last := generatorget;
       Take := generatortake;
       SetInterval := generatorsetinterval;
       PipelineInterval := generatorpipelineinterval;
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}

func GeneratorSetup() {
	log.Debug("Running Generator setup")
	GoStructRegistry.RegisterUserdef(&RegisteredType{GenDefMap: true, Factory: func(env *Zlisp, h *SexpHash) (interface{}, error) {
		return &GEN{HistorySize: int64(DefaultHistorySize)}, nil
	}}, true, "GEN")
}
