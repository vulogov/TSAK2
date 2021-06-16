package dsl

import (
	"fmt"

	"github.com/edwingeng/deque"
	"github.com/glycerine/zygomys/zygo"
	. "github.com/glycerine/zygomys/zygo"
	"github.com/pieterclaerhout/go-log"

	"github.com/lrita/cmap"
)

type TelemetryObservationMatrix struct {
	Matrix               cmap.Cmap
	Defaults             cmap.Cmap
	ObservabilityHorizon int64 `json:"ObservabilityHorizon" msg:"ObservabilityHorizon"`
}

var DefaultObservabilityHorizon = 64
var TOM = TelemetryObservationMatrix{ObservabilityHorizon: int64(DefaultObservabilityHorizon)}

func makekey(source, name string) string {
	return fmt.Sprintf("%s@%s", source, name)
}

func (tom *TelemetryObservationMatrix) _create(source, name string, defval interface{}) bool {
	log.Debugf("Creating TDATA for %s@%s %v", source, name, tom.ObservabilityHorizon)
	key := makekey(source, name)
	_, ok := tom.Defaults.Load(key)
	if ok {
		return false
	}
	tom.Defaults.Store(key, defval)
	r := deque.NewDeque()
	tom.Matrix.Store(key, r)
	return true
}

func (tom *TelemetryObservationMatrix) CreateFloat(source, name string, defval float64) bool {
	return tom._create(source, name, defval)
}

func (tom *TelemetryObservationMatrix) CreateInt(source, name string, defval int64) bool {
	return tom._create(source, name, defval)
}

func (tom *TelemetryObservationMatrix) CreateString(source, name string, defval string) bool {
	return tom._create(source, name, defval)
}

func (tom *TelemetryObservationMatrix) _add(source, name string, val interface{}) int {
	key := makekey(source, name)
	_m, ok := tom.Matrix.Load(key)
	if !ok {
		tom._create(source, name, val)
		_m, ok = tom.Matrix.Load(key)
	}
	m := _m.(deque.Deque)
	if m.Len() > int(tom.ObservabilityHorizon) {
		m.Dequeue()
	}
	m.Enqueue(val)
	return m.Len()
}

func (tom *TelemetryObservationMatrix) Horizon() int64 {
	return tom.ObservabilityHorizon
}

func (tom *TelemetryObservationMatrix) AddFloat(source, name string, val float64) int {
	return tom._add(source, name, val)
}

func (tom *TelemetryObservationMatrix) AddInt(source, name string, val int64) int {
	return tom._add(source, name, val)
}

func (tom *TelemetryObservationMatrix) AddString(source, name string, val string) int {
	return tom._add(source, name, val)
}

func (tom *TelemetryObservationMatrix) Get(source, name string) interface{} {
	key := makekey(source, name)
	_m, ok := tom.Matrix.Load(key)
	if !ok {
		return nil
	}
	m := _m.(deque.Deque)
	if m.Len() == 0 {
		d, ok := tom.Defaults.Load(key)
		if !ok {
			return nil
		}
		return d
	}
	return m.Dequeue()
}

func (tom *TelemetryObservationMatrix) Len(source, name string) int {
	key := makekey(source, name)
	_m, ok := tom.Matrix.Load(key)
	if !ok {
		return 0
	}
	m := _m.(deque.Deque)
	return m.Len()
}

func (tom *TelemetryObservationMatrix) _setDefault(source, name string, defval interface{}) bool {
	key := makekey(source, name)
	_, ok := tom.Defaults.Load(key)
	if !ok {
		return false
	}
	tom.Defaults.Store(key, defval)
	return true
}

func (tom *TelemetryObservationMatrix) SetDefaultFloat(source, name string, defval float64) bool {
	return tom._setDefault(source, name, defval)
}

func (tom *TelemetryObservationMatrix) SetDefaultInt(source, name string, defval int64) bool {
	return tom._setDefault(source, name, defval)
}

func (tom *TelemetryObservationMatrix) SetDefaultString(source, name string, defval string) bool {
	return tom._setDefault(source, name, defval)
}

func (tom *TelemetryObservationMatrix) Sample(source, name string) (res []interface{}) {
	key := makekey(source, name)
	_m, ok := tom.Matrix.Load(key)
	if !ok {
		return
	}
	m := _m.(deque.Deque)
	n := m.Len()
	res = make([]interface{}, n)
	for c := 0; c < n; c++ {
		res[c] = m.Peek(c)
	}
	return
}

func (tom *TelemetryObservationMatrix) _matrix(source, name []string, defval interface{}) int {
	var n = 0
	for _, i := range source {
		for _, j := range name {
			if tom._create(i, j, defval) {
				n += 1
			}
		}
	}
	return n
}

func (tom *TelemetryObservationMatrix) Float(source, name []string, val float64) int {
	return tom._matrix(source, name, val)
}
func (tom *TelemetryObservationMatrix) Int(source, name []string, val int64) int {
	return tom._matrix(source, name, val)
}
func (tom *TelemetryObservationMatrix) String(source, name []string, val string) int {
	return tom._matrix(source, name, val)
}

func TOMMake(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) != 3 {
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

	switch v := args[2].(type) {
	case *SexpFloat:
		n += TOM.Float(ar1, ar2, v.Val)
	case *SexpInt:
		n += TOM.Int(ar1, ar2, v.Val)
	case *SexpStr:
		n += TOM.String(ar1, ar2, v.S)
	default:
		n += TOM.String(ar1, ar2, args[2].SexpString(nil))
	}

	return &SexpInt{Val: int64(n)}, nil
}

func TOMAdd(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) != 3 {
		return SexpNull, WrongNargs
	}
	if !IsString(args[0]) {
		return SexpNull, fmt.Errorf("First argument must be string")
	}
	if !IsString(args[1]) {
		return SexpNull, fmt.Errorf("Second argument must be string")
	}

	n := 0

	switch v := args[2].(type) {
	case *SexpInt:
		n += TOM.AddInt(AsString(args[0]), AsString(args[1]), v.Val)
	case *SexpFloat:
		n += TOM.AddFloat(AsString(args[0]), AsString(args[1]), v.Val)
	case *SexpStr:
		n += TOM.AddString(AsString(args[0]), AsString(args[1]), AsString(args[2]))
	default:
		n += TOM.AddString(AsString(args[0]), AsString(args[1]), AsString(args[2]))
	}
	return &SexpInt{Val: int64(n)}, nil
}

func TOMGet(env *Zlisp, name string, args []Sexp) (Sexp, error) {

	if len(args) != 2 {
		return SexpNull, WrongNargs
	}
	if !IsString(args[0]) {
		return SexpNull, fmt.Errorf("First argument must be string")
	}
	if !IsString(args[1]) {
		return SexpNull, fmt.Errorf("Second argument must be string")
	}

	v := TOM.Get(AsString(args[0]), AsString(args[1]))
	switch e := v.(type) {
	case int:
	case int64:
	case uint32:
	case uint64:
		return &SexpInt{Val: int64(e)}, nil
	case float64:
		return &SexpFloat{Val: float64(e)}, nil
	case string:
		return &SexpStr{S: e}, nil
	}
	return SexpNull, fmt.Errorf("TOM returned unexpected telemetry type")
}

func TelemetryObservationFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"tommake": TOMMake,
		"tomadd":  TOMAdd,
		"tomget":  TOMGet,
	}
}

func TelemetryObservationPackageSetup(cfg *zygo.ZlispConfig, env *zygo.Zlisp) {
	myPkg := `(def observation (package "observation"
     { Make := tommake;
       Add := tomadd;
       Get := tomget;
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}

func TelemetryObservationMatrixSetup() {
	log.Debug("Running TelemetryObservationMatrix setup")
	GoStructRegistry.RegisterUserdef(&RegisteredType{GenDefMap: true, Factory: func(env *Zlisp, h *SexpHash) (interface{}, error) {
		return &TelemetryObservationMatrix{ObservabilityHorizon: int64(DefaultObservabilityHorizon)}, nil
	}}, true, "TelemetryObservationMatrix")
}
