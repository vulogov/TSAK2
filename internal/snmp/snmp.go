package snmp

import (
	"log"
	"os"
	"time"

	// "github.com/pieterclaerhout/go-log"
	"github.com/twsnmp/gosnmp"

	"github.com/vulogov/TSAK2/internal/conf"
	"github.com/vulogov/TSAK2/internal/signal"
)

var Gsnmp = &gosnmp.GoSNMP{}
var AgentSnmp = &gosnmp.GoSNMPAgent{}

func SnmpAgentLoop() {
	// log.Debug("Entering SNMP Agent service loop")
	AgentSnmp.Start()
	for !signal.ExitRequested() {
		time.Sleep(200 * time.Millisecond)
	}
	// log.Debug("Terminating internal SNMP agent")
	AgentSnmp.Stop()
}

func InitSNMPAgent() {
	// log.Debug("Starting internal SNMP Agent")
	Gsnmp.Target = *conf.SNMPListen
	Gsnmp.Port = uint16(*conf.SNMPPort)
	Gsnmp.Community = *conf.SNMPCommunity
	Gsnmp.Version = gosnmp.Version2c
	Gsnmp.Timeout = time.Duration(time.Second * 3)
	Gsnmp.Retries = 0
	Gsnmp.Logger = log.New(os.Stdout, "", 0)
	AgentSnmp.Port = int(*conf.SNMPPort)
	AgentSnmp.IPAddr = *conf.SNMPListen
	AgentSnmp.Snmp = Gsnmp
	AgentSnmp.SupportSnmpMIB = false
	AgentSnmp.Logger = log.New(os.Stdout, "", 0)
	go SnmpAgentLoop()
}
