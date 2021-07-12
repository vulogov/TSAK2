package dsl

import (
	"fmt"
	"reflect"
	"time"

	. "github.com/glycerine/zygomys/zygo"
	"github.com/pieterclaerhout/go-log"
	"github.com/twsnmp/gosnmp"

	"github.com/vulogov/TSAK2/internal/conf"
	"github.com/vulogov/TSAK2/internal/snmp"
)

func tomSnmpGet(oid string) interface{} {
	log.Debugf("Requested: %v", oid)
	res := TOM.Get("snmp", oid)
	return res
}

func snmpAgentSet(oid string, defval interface{}) {
	switch e := defval.(type) {
	case int, int64, uint32, uint64:
		TOM.AddInt("snmp", oid, int64(e.(int64)))
		TOM.SetSize("snmp", oid, 34)
		log.Debugf("Configured: %v", oid)
		snmp.AgentSnmp.AddMibList(oid, gosnmp.Integer, func(oid string) interface{} { return TOM.Get("snmp", oid) })
	case float64:
		TOM.AddFloat("snmp", oid, float64(e))
		snmp.AgentSnmp.AddMibList(oid, gosnmp.OpaqueFloat, func(oid string) interface{} { return TOM.Get("snmp", oid) })
	case string:
		TOM.AddString("snmp", oid, string(e))
		snmp.AgentSnmp.AddMibList(oid, gosnmp.OctetString, func(oid string) interface{} { return TOM.Get("snmp", oid) })
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

func SnmpTrapsRecvd(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) != 0 {
		return SexpNull, WrongNargs
	}
	return &SexpInt{Val: int64(snmp.TrapData.Len())}, nil
}

func SnmpTrapsRecv(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) != 0 {
		return SexpNull, WrongNargs
	}
	res, err := MakeHash(make([]Sexp, 0), "hash", env)
	PanicOn(err)
	if snmp.TrapData.Len() > 0 {
		d := snmp.TrapData.PopFront()
		id := d.(snmp.TRAPData).ID
		key := d.(snmp.TRAPData).Key
		val := d.(snmp.TRAPData).Val
		res.HashSet(&SexpStr{S: "ID"}, &SexpInt{Val: int64(id)})
		res.HashSet(&SexpStr{S: "OID"}, &SexpStr{S: string(key)})
		switch e := val.(type) {
		case int:
			res.HashSet(&SexpStr{S: "Value"}, &SexpInt{Val: int64(e)})
		case int64:
			res.HashSet(&SexpStr{S: "Value"}, &SexpInt{Val: int64(e)})
		case uint64:
			res.HashSet(&SexpStr{S: "Value"}, &SexpInt{Val: int64(e)})
		case uint32:
			res.HashSet(&SexpStr{S: "Value"}, &SexpInt{Val: int64(e)})
		case int32:
			res.HashSet(&SexpStr{S: "Value"}, &SexpInt{Val: int64(e)})
		case float64:
			res.HashSet(&SexpStr{S: "Value"}, &SexpFloat{Val: float64(e)})
		case string:
			res.HashSet(&SexpStr{S: "Value"}, &SexpStr{S: string(e)})
		default:
			log.Debugf("Unknown trap packet: %v %v", d.(snmp.TRAPData), reflect.TypeOf(val))
		}
	}
	return res, nil
}

func SnmpTrapsSnd(env *Zlisp, name string, args []Sexp) (Sexp, error) {
	if len(args) != 3 {
		return SexpNull, WrongNargs
	}
	if !IsString(args[0]) {
		return SexpNull, fmt.Errorf("First argument must be string")
	}
	if !IsInt(args[1]) {
		return SexpNull, fmt.Errorf("Second argument must be integer")
	}
	if !IsHash(args[2]) {
		return SexpNull, fmt.Errorf("Third argument must be hash")
	}
	gsnmp := &gosnmp.GoSNMP{
		Target:    AsString(args[0]),
		Port:      uint16(AsAny(args[1]).(int64)),
		Community: *conf.SNMPCommunity,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Retries:   1,
		MaxOids:   gosnmp.MaxOids,
	}
	err := gsnmp.Connect()
	PanicOn(err)
	defer gsnmp.Conn.Close()
	n := 0
	trap := gosnmp.SnmpTrap{
		Variables: make([]gosnmp.SnmpPDU, 0),
	}
	switch e := args[2].(type) {
	case *SexpHash:
		for _, v := range e.Map {
			for _, p := range v {
				fmt.Println(p.Head, p.Tail)
				pdu := gosnmp.SnmpPDU{
					Name: AsString(p.Head),
				}
				switch e1 := p.Tail.(type) {
				case *SexpInt:
					pdu.Type = gosnmp.Integer
					pdu.Value = int(e1.Val)
				case *SexpStr:
					pdu.Type = gosnmp.OctetString
					pdu.Value = string(e1.S)
				case *SexpFloat:
					pdu.Type = gosnmp.OpaqueFloat
					pdu.Value = float32(e1.Val)
				default:
					continue
				}
				trap.Variables = append(trap.Variables, pdu)
			}
		}
	}
	_, err = gsnmp.SendTrap(trap)
	PanicOn(err)
	return &SexpInt{Val: int64(n)}, nil
}

func SnmpMetricFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{
		"snmpmetricset": SnmpAgentSet,
		"snmptraprecvd": SnmpTrapsRecvd,
		"snmptraprecv":  SnmpTrapsRecv,
		"snmptrapsnd":   SnmpTrapsSnd,
	}
}

func SnmpMetricPackageSetup(cfg *ZlispConfig, env *Zlisp) {
	myPkg := `(def snmpmetric (package "snmpmetric"
     { Make := (fn [oids defval] (observation.Make ["snmp"] oids defval));
       Add := (fn [oid val] (observation.Add "snmp" oid val));
       Get := (fn [oid] (observation.Get "snmp" oid ));
			 Has := (fn [oid] (observation.Has "snmp" oid ));
			 Set := snmpmetricset ;
			 Traps := snmptraprecvd ;
			 Trap := snmptraprecv ;
			 TrapSend := snmptrapsnd ;
     }
  ))`
	_, err := env.EvalString(myPkg)
	PanicOn(err)
}
