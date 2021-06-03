package conf

import (
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/rs/xid"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	seed    = time.Now().UTC().UnixNano()
	NG      = namegenerator.NewNameGenerator(seed)
	App     = kingpin.New("tsak2", "[ tsak2 ] Telemetry Swiss Army Knife")
	Debug   = App.Flag("debug", "Enable debug mode.").Default("false").Bool()
	Color   = App.Flag("color", "--color : Enable colors on terminal --no-color : Disable colors .").Default("true").Bool()
	ID      = App.Flag("id", "Unique application ID").Default(xid.New().String()).String()
	Name    = App.Flag("name", "Application name").Default(NG.Generate()).String()
	VBanner = App.Flag("banner", "Display [ tsak2 ] banner .").Default("false").Bool()

	// SNMP-related configuration
	SNMPCommunity = App.Flag("community", "SNMP 2c community string").Default("public").String()
	SNMPUser      = App.Flag("snmpuser", "SNMP 3 username").Default("v3user").String()
	SNMPAuthPass  = App.Flag("snmpauth", "SNMP 3 authentication passphrase").Default("v3password").String()
	SNMPAuthPriv  = App.Flag("snmppriv", "SNMP 3 privacy passphrase").Default("v3priv").String()
	SNMPListen    = App.Flag("snmplisten", "Address for internal TSAK SNMP agent").Default("127.0.0.1:6161").String()

	// Bootstrap-related
	BootStrap = App.Flag("boot", "TSAK script for the environment bootstrap").ExistingFile()

	Version = App.Command("version", "Display information about [ tsak2 ]")
	VTable  = Version.Flag("table", "Display [ tsak2 ] inner information .").Default("true").Bool()

	Shell = App.Command("shell", "Run tsak2 in interactive shell")
	Run   = App.Command("run", "Run tsak2 in non-interactive mode")
)
