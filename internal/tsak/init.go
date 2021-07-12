package tsak

import (
	"github.com/pieterclaerhout/go-log"

	tlog "github.com/vulogov/TSAK2/internal/log"
	"github.com/vulogov/TSAK2/internal/signal"
	"github.com/vulogov/TSAK2/internal/snmp"
)

func Init() {
	tlog.Init()
	log.Debug("[ TSAK2 ] tsak.Init() is reached")
	signal.InitSignal()
	snmp.InitSNMPAgent()
	snmp.InitSNMPTrapReceiver()
}
