package dsl

import (
	"fmt"

	ring "github.com/zealws/golang-ring"

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

func (tom *TelemetryObservationMatrix) Create(source, name string, defval interface{}) bool {
	key := makekey(source, name)
	_, ok := tom.Defaults.Load(key)
	if ok {
		return false
	}
	tom.Defaults.Store(key, defval)
	r := ring.Ring{}
	r.SetCapacity(int(tom.ObservabilityHorizon))
	tom.Matrix.Store(key, r)
	return true
}

func (tom *TelemetryObservationMatrix) Add(source, name string, val interface{}) int {
	key := makekey(source, name)
	m, ok := tom.Matrix.Load(key)
	if !ok {
		tom.Create(source, name, val)
		m, ok = tom.Matrix.Load(key)
	}
	m.(*ring.Ring).Enqueue(val)
	return m.(*ring.Ring).ContentSize()
}

func (tom *TelemetryObservationMatrix) Get(source, name string) interface{} {
	key := makekey(source, name)
	m, ok := tom.Matrix.Load(key)
	if !ok {
		return nil
	}
	if m.(*ring.Ring).ContentSize() == 0 {
		d, ok := tom.Defaults.Load(key)
		if !ok {
			return nil
		}
		return d
	}
	return m.(*ring.Ring).Dequeue()
}

func (tom *TelemetryObservationMatrix) Len(source, name string) int {
	key := makekey(source, name)
	m, ok := tom.Matrix.Load(key)
	if !ok {
		return 0
	}
	return m.(*ring.Ring).ContentSize()
}

func (tom *TelemetryObservationMatrix) SetDefault(source, name string, defval interface{}) bool {
	key := makekey(source, name)
	_, ok := tom.Defaults.Load(key)
	if !ok {
		return false
	}
	tom.Defaults.Store(key, defval)
	return true
}
