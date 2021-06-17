package dsl

import (
	"fmt"

	. "github.com/glycerine/zygomys/zygo"
	"github.com/twsnmp/gosnmp"

	"github.com/vulogov/TSAK2/internal/snmp"
)

func tomSnmpGet(oid string) interface{} {
	res := TOM.Get("snmp", oid)
	return res
}

func snmpAgentSet(oid string, defval interface{}) {
	switch e := defval.(type) {
	case int, int64, uint32, uint64:
		TOM.AddInt("snmp", oid, int64(e.(int64)))
		TOM.SetSize("snmp", oid, 34)
		snmp.AgentSnmp.AddMibList(oid, gosnmp.Integer, tomSnmpGet)
	case float64:
		TOM.AddFloat("snmp", oid, float64(e))
		snmp.AgentSnmp.AddMibList(oid, gosnmp.OpaqueFloat, tomSnmpGet)
	case string:
		TOM.AddString("snmp", oid, string(e))
		snmp.AgentSnmp.AddMibList(oid, gosnmp.OctetString, tomSnmpGet)
	}
}

func SnmpAgentSet(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) != 2 {
		return SexpNull, WrongNargs
	}
	n := 0
	if IsArray(args[0]) {
		oids := ArrayofStringsToArray(args[0])
		for _, o := range oids {
			snmpAgentSet(o, AsAny(args[1]))
			n += 1
		}
	} else if IsString(args[0]) {
		snmpAgentSet(AsString(args[0]), AsAny(args[1]))
		n += 1
	} else {
		return SexpNull, fmt.Errorf("First argument must be ether string or array")
	}
	return &SexpInt{Val: int64(n)}, nil
}

func SnmpMetricFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"snmpmetricset": SnmpAgentSet,
	}
}

func SnmpMetricPackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def snmpmetric (package "snmpmetric"
     { Make := (fn [oids defval] (observation.Make ["snmp"] oids defval));
       Add := (fn [oid val] (observation.Add "snmp" oid val));
       Get := (fn [oid] (observation.Get "snmp" oid ));
			 Has := (fn [oid] (observation.Has "snmp" oid ));
			 Set := snmpmetricset ;
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}
