package conf

import (
	"fmt"
	"os"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/rs/xid"
	"gopkg.in/alecthomas/kingpin.v2"
)

type filelist []string

func (i *filelist) Set(value string) error {
	_, err := os.Stat(value)
	if os.IsNotExist(err) {
		return fmt.Errorf("Script file '%s' not found", value)
	} else {
		*i = append(*i, value)
		return nil
	}
}

func (i *filelist) String() string {
	return ""
}

func (i *filelist) IsCumulative() bool {
	return true
}

func FileList(s kingpin.Settings) (target *[]string) {
	target = new([]string)
	s.SetValue((*filelist)(target))
	return
}

var (
	seed    = time.Now().UTC().UnixNano()
	NG      = namegenerator.NewNameGenerator(seed)
	App     = kingpin.New("TSAK2", "[ TSAK2 ] Telemetry Swiss Army Knife")
	Debug   = App.Flag("debug", "Enable debug mode.").Default("false").Bool()
	Color   = App.Flag("color", "--color : Enable colors on terminal --no-color : Disable colors .").Default("true").Bool()
	ID      = App.Flag("id", "Unique application ID").Default(xid.New().String()).String()
	Name    = App.Flag("name", "Application name").Default(NG.Generate()).String()
	VBanner = App.Flag("banner", "Display [ TSAK2 ] banner .").Default("false").Bool()

	// SNMP-related configuration
	SNMPCommunity = App.Flag("community", "SNMP 2c community string").Default("public").String()
	SNMPListen    = App.Flag("snmplisten", "IP Address for internal TSAK SNMP agent").Default("127.0.0.1").String()
	SNMPPort      = App.Flag("snmpport", "Port for internal TSAK SNMP agent").Default("6161").Int()
	SNMPMibsdb    = App.Flag("mibs", "Path to SNMP MIB database").Default("./mibs/nri-snmp.db").ExistingFile()
	TRAPPort      = App.Flag("trapport", "Port for internal TSAK TRAP receiver").Default("6162").Int()
	TRAPListen    = App.Flag("traplisten", "IP Address for internal TSAK TRAP Receiver").Default("127.0.0.1").String()

	// Bootstrap-related
	BootStrap = App.Flag("boot", "TSAK script for the environment bootstrap").ExistingFile()

	// Generator and pipeline
	GenInterval = App.Flag("gen", "How often to run data generator").Default("60").Int()

	// New Relic related
	NRapi     = App.Flag("nrapi", "New Relic API key").Default("").String()
	NRapiq    = App.Flag("nrapiq", "New Relic Query API key").Default("").String()
	NRaccount = App.Flag("nraccount", "New Relic Account ID").Default("").String()
	NRlogapi  = App.Flag("nrlogapi", "New Relic Log API").Default("https://log-api.newrelic.com/log/v1").String()
	NRevtapi  = App.Flag("nrevtapi", "New Relic Event API").Default("https://insights-collector.newrelic.com/v1/accounts/%s/events").String()
	NRmetapi  = App.Flag("nrmetapi", "New Relic Metrics API").Default("https://metric-api.newrelic.com/metric/v1").String()
	NRqryapi  = App.Flag("nrqryapi", "New Relic Query API").Default("https://insights-api.newrelic.com/v1/accounts/%s/query?nrql=%s").String()

	Version = App.Command("version", "Display information about [ TSAK2 ]")
	VTable  = Version.Flag("table", "Display [ TSAK2 ] inner information .").Default("true").Bool()

	Shell      = App.Command("shell", "Run TSAK2 in interactive shell")
	Run        = App.Command("run", "Run TSAK2 in non-interactive mode")
	ShowResult = Run.Flag("result", "Display result of scripts execution as it returned by LISP").Default("true").Bool()
	ERloop     = Run.Flag("erloop", "ExitRequest event loop").Default("false").Bool()
	WGloop     = Run.Flag("loop", "WorkGroup event loop").Default("false").Bool()
	Scripts    = FileList(Run.Arg("Scripts", "TSAK-scripts to execute"))
)
