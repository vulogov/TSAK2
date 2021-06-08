package dsl

import "github.com/glycerine/zygomys/zygo"

func SnmpMetricPackageSetup(cfg *zygo.ZlispConfig, env *zygo.Zlisp) {
	myPkg := `(def snmpmetric (package "snmpmetric"
     { Make := (fn [oids defval] (observation.Make ["snmp"] oids defval));
       Add := (fn [oid val] (observation.Add "snmp" oid val));
       Get := (fn [oid] (observation.Get "snmp" oid ));
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}
