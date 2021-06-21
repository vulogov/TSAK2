package dsl

import (
	"database/sql"
	"os"

	. "github.com/glycerine/zygomys/zygo"
	"github.com/lrita/cmap"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pieterclaerhout/go-log"
	"github.com/twsnmp/gosnmp"

	"github.com/vulogov/TSAK2/internal/conf"
	"github.com/vulogov/TSAK2/internal/snmp"
)

type MIBS struct {
	DBname  string
	TDB     *sql.DB
	MibsMap cmap.Cmap
	TiMap   cmap.Cmap
	TsMap   cmap.Cmap
}

func (mibs *MIBS) Open(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	mibs.DBname = name
	mibs.TDB, err = sql.Open("sqlite3", name)
	if err != nil {
		log.Errorf("Fail to open %s: %v", name, err)
		return false
	}
	return true
}

func (mibs *MIBS) Init() bool {
	return mibs.Open(*conf.SNMPMibsdb)
}

func (mibs *MIBS) IsOpen() bool {
	if mibs.TDB == nil {
		return false
	}
	return true
}

func (mibs *MIBS) AddMib(name string) int {
	if !mibs.IsOpen() {
		log.Errorf("MIB db %s is not opened", mibs.DBname)
		return 0
	}
	stmt, err := mibs.TDB.Prepare("select ID from MIB where NAME=?")
	if err != nil {
		log.Errorf("Error in preparation MIB query: %v", err)
		return 0
	}
	rows, err := stmt.Query(name)
	if err != nil {
		log.Errorf("Error executing MIB query: %v", err)
		return 0
	}
	for rows.Next() {
		var n int
		rows.Scan(&n)
		mibs.MibsMap.Store(name, n)
		mstmt, err := mibs.TDB.Prepare("select KEY,TI,TS from OBJECT where MIB=?")
		if err != nil {
			log.Errorf("Error in preparation OBJECT query: %v", err)
			return 0
		}
		mrows, err := mstmt.Query(n)
		if err != nil {
			log.Errorf("Error executing OBJECT query: %v", err)
			return 0
		}
		for mrows.Next() {
			var k, ts string
			var ti int
			mrows.Scan(&k, &ti, &ts)
			mibs.TiMap.LoadOrStore(k, ti)
			mibs.TsMap.LoadOrStore(k, ts)
		}
		return n
	}
	log.Errorf("MIB query returns no value")
	return 0
}

func (mibs *MIBS) AllMibs() []string {
	var res []string
	if !mibs.IsOpen() {
		log.Errorf("MIB db %s is not opened", mibs.DBname)
		return res
	}
	res = make([]string, 0)
	stmt, err := mibs.TDB.Prepare("select NAME from MIB")
	if err != nil {
		log.Errorf("Error in preparation MIB query: %v", err)
		return res
	}
	rows, err := stmt.Query()
	if err != nil {
		log.Errorf("Error executing MIB query: %v", err)
		return res
	}
	for rows.Next() {
		var n string
		rows.Scan(&n)
		res = append(res, n)
	}
	return res
}

func (mibs *MIBS) Mibs() []string {
	var res []string
	if !mibs.IsOpen() {
		log.Errorf("MIB db %s is not opened", mibs.DBname)
		return res
	}
	res = make([]string, 0)
	mibs.MibsMap.Range(func(ki, vi interface{}) bool {
		k, _ := ki.(string), vi.(int)
		res = append(res, k)
		return true
	})
	return res
}

func (mibs *MIBS) Oids() []string {
	var res []string
	if !mibs.IsOpen() {
		log.Errorf("MIB db %s is not opened", mibs.DBname)
		return res
	}
	res = make([]string, 0)
	mibs.TiMap.Range(func(ki, vi interface{}) bool {
		k, _ := ki.(string), vi.(int)
		res = append(res, k)
		return true
	})
	return res
}

func (mibs *MIBS) Ti(oid string) int {
	if !mibs.IsOpen() {
		log.Errorf("MIB db %s is not opened", mibs.DBname)
		return 0
	}
	d, ok := mibs.TiMap.Load(oid)
	if !ok {
		return 0
	}
	return d.(int)
}

func (mibs *MIBS) Ts(oid string) string {
	if !mibs.IsOpen() {
		log.Errorf("MIB db %s is not opened", mibs.DBname)
		return ""
	}
	d, ok := mibs.TsMap.Load(oid)
	if !ok {
		return ""
	}
	return d.(string)
}

func (mibs *MIBS) Export() int {
	n := 0
	mibs.TiMap.Range(func(ki, vi interface{}) bool {
		k, v := ki.(string), vi.(int)
		if !TOM.Has("snmp", k) {
			tsv, ok := mibs.TsMap.Load(k)
			if !ok {
				tsv = ""
			}
			switch int(v) {
			case 3, 4, 64, 18, 6, 19, 20, 21, 26, 7:
				if tsv == "IpAddress" {
					TOM.AddString("snmp", k, "127.0.0.1")
					snmp.AgentSnmp.AddMibList(k, gosnmp.IPAddress, tomSnmpGet)
					log.Debugf("Exporting %s as IPAddr", k)
				} else if tsv == "Bits" {
					TOM.AddString("snmp", k, "")
					snmp.AgentSnmp.AddMibList(k, gosnmp.BitString, tomSnmpGet)
					log.Debugf("Exporting %s as BitString", k)
				} else if tsv == "OBJECT-IDENTIFIER" {
					TOM.AddString("snmp", k, k)
					snmp.AgentSnmp.AddMibList(k, gosnmp.ObjectIdentifier, tomSnmpGet)
					log.Debugf("Exporting %s as ObjectIdentifier", k)
				} else if tsv == "ObjectDescriptor" {
					TOM.AddString("snmp", k, "")
					snmp.AgentSnmp.AddMibList(k, gosnmp.ObjectDescription, tomSnmpGet)
					log.Debugf("Exporting %s as ObjectDescription", k)
				} else {
					TOM.AddString("snmp", k, "")
					snmp.AgentSnmp.AddMibList(k, gosnmp.OctetString, tomSnmpGet)
					log.Debugf("Exporting %s as OctetString", k)
				}
				n += 1
			case 2, 0, 23:
				if tsv == "Counter" {
					TOM.AddInt("snmp", k, 0)
					TOM.SetSize("snmp", k, 33)
					snmp.AgentSnmp.AddMibList(k, gosnmp.Counter32, tomSnmpGet)
					log.Debugf("Exporting %s as Counter", k)
				} else if tsv == "Counter32" {
					TOM.AddInt("snmp", k, 0)
					TOM.SetSize("snmp", k, 33)
					snmp.AgentSnmp.AddMibList(k, gosnmp.Counter32, tomSnmpGet)
					log.Debugf("Exporting %s as Counter32", k)
				} else if tsv == "Counter64" {
					TOM.AddInt("snmp", k, 0)
					TOM.SetSize("snmp", k, 65)
					snmp.AgentSnmp.AddMibList(k, gosnmp.Counter64, tomSnmpGet)
					log.Debugf("Exporting %s as Counter64", k)
				} else if tsv == "Gauge" {
					TOM.AddInt("snmp", k, 0)
					TOM.SetSize("snmp", k, 33)
					snmp.AgentSnmp.AddMibList(k, gosnmp.Gauge32, tomSnmpGet)
					log.Debugf("Exporting %s as Gauge", k)
				} else if tsv == "Gauge32" {
					TOM.AddInt("snmp", k, 0)
					TOM.SetSize("snmp", k, 33)
					snmp.AgentSnmp.AddMibList(k, gosnmp.Gauge32, tomSnmpGet)
					log.Debugf("Exporting %s as Gauge32", k)
				} else if tsv == "TimeTicks" {
					log.Debugf("TimeTicks type for %s not yet supported", k)
				} else {
					TOM.AddInt("snmp", k, 0)
					TOM.SetSize("snmp", k, 34)
					snmp.AgentSnmp.AddMibList(k, gosnmp.Integer, tomSnmpGet)
					log.Debugf("Exporting %s as Integer", k)
				}
				n += 1
			case 9:
				TOM.AddFloat("snmp", k, float64(0.0))
				snmp.AgentSnmp.AddMibList(k, gosnmp.OpaqueFloat, tomSnmpGet)
				log.Debugf("Exporting %s as Real", k)
				n += 1
			default:
				log.Debugf("Not exporting %s as ti: %d ts: %s", k, v, tsv)
			}
		}
		return true
	})
	return n
}

func MIBSFunctions() map[string]ZlispUserFunction {
	return map[string]ZlispUserFunction{}
}

func MIBSDatatypeSetup() {
	log.Debug("Running MIBS setup")
	GoStructRegistry.RegisterUserdef(&RegisteredType{GenDefMap: true, Factory: func(env *Zlisp, h *SexpHash) (interface{}, error) {
		return &MIBS{DBname: *conf.SNMPMibsdb, TDB: nil}, nil
	}}, true, "MIBS")
}
