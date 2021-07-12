package snmp

import (
	"fmt"
	// "log"
	"net"
	// "os"
	"time"

	"github.com/edwingeng/deque"

	// "github.com/pieterclaerhout/go-log"
	"github.com/twsnmp/gosnmp"

	"github.com/vulogov/TSAK2/internal/conf"
	"github.com/vulogov/TSAK2/internal/signal"
)

var Gsnmp = &gosnmp.GoSNMP{}
var GsnmpTrap = &gosnmp.GoSNMP{}

var AgentSnmp = &gosnmp.GoSNMPAgent{}
var TrapListen = gosnmp.NewTrapListener()
var TrapData = deque.NewDeque()

type TRAPData struct {
	ID  uint32
	Key string
	Val interface{}
}

func trapHandler(packet *gosnmp.SnmpPacket, addr *net.UDPAddr) {
	for _, v := range packet.Variables {
		p := TRAPData{ID: packet.RequestID, Key: v.Name, Val: v.Value}
		// log.Debugf("TRAP receiver: %v", p)
		TrapData.PushBack(p)
	}
}

func SnmpAgentLoop() {
	// log.Debug("Entering SNMP Agent service loop")
	AgentSnmp.Start()
	for !signal.ExitRequested() {
		time.Sleep(200 * time.Millisecond)
	}
	// log.Debug("Terminating internal SNMP agent")
	AgentSnmp.Stop()
}

func TrapReceiverLoop() {
	// log.Debug("Entering TRAP receiver service loop")
	defer TrapListen.Close()
	err := TrapListen.Listen(net.JoinHostPort(*conf.TRAPListen, fmt.Sprintf("%d", *conf.TRAPPort)))
	if err != nil {
		// log.Errorf("TRAP listener error: %v", err)
		signal.ExitRequest()
	}
	// log.Debug("Exiting TRAP receiver service loop")
}

func InitSNMPAgent() {
	// log.Debug("Starting internal SNMP Agent")
	Gsnmp.Target = *conf.SNMPListen
	Gsnmp.Port = uint16(*conf.SNMPPort)
	Gsnmp.Community = *conf.SNMPCommunity
	Gsnmp.Version = gosnmp.Version2c
	Gsnmp.Timeout = time.Duration(time.Second * 3)
	Gsnmp.Retries = 1
	// Gsnmp.Logger = gosnmp.NewLogger(log.New(os.Stdout, "", 0))
	AgentSnmp.Port = int(*conf.SNMPPort)
	AgentSnmp.IPAddr = *conf.SNMPListen
	AgentSnmp.Snmp = Gsnmp
	AgentSnmp.SupportSnmpMIB = false
	// AgentSnmp.Logger = gosnmp.NewLogger(log.New(os.Stdout, "", 0))
	go SnmpAgentLoop()
}

func InitSNMPTrapReceiver() {
	// log.Debug("Starting internal Trap receiver")
	GsnmpTrap.Target = *conf.TRAPListen
	GsnmpTrap.Port = uint16(*conf.TRAPPort)
	GsnmpTrap.Community = *conf.SNMPCommunity
	GsnmpTrap.Version = gosnmp.Version2c
	GsnmpTrap.Timeout = time.Duration(time.Second * 3)
	GsnmpTrap.Retries = 0
	TrapListen.Params = GsnmpTrap
	TrapListen.OnNewTrap = trapHandler
	go TrapReceiverLoop()
}
