package snmp

import (
	// "github.com/vulogov/TSAK2/internal/dsl"
	"github.com/pieterclaerhout/go-log"
)

func SnmpAgentLoop() {
	log.Debug("Entering SNMP Agent service loop")

}

func InitSNMPAgent() {
	log.Debug("Starting internal SNMP Agent")
	go SnmpAgentLoop()
}
