package banner

import (
	"os"

	"github.com/mgutz/ansi"
	"github.com/tomlazar/table"

	"github.com/vulogov/TSAK2/internal/conf"
)

func Table() {
	var cfg table.Config

	if !*conf.VTable {
		return
	}

	cfg.ShowIndex = true
	if *conf.Color {
		cfg.Color = true
		cfg.AlternateColors = true
		cfg.TitleColorCode = ansi.ColorCode("white+buf")
		cfg.AltColorCodes = []string{"", ansi.ColorCode("white:grey+h")}
	} else {
		cfg.Color = false
		cfg.AlternateColors = false
		cfg.TitleColorCode = ansi.ColorCode("white+buf")
		cfg.AltColorCodes = []string{"", ansi.ColorCode("white:grey+h")}
	}
	if *conf.VTable {
		tab := table.Table{
			Headers: []string{"Description", "Value"},
			Rows: [][]string{
				{"Version", conf.BVersion},
				{"Application ID", *conf.ID},
				{"Application name", *conf.Name},
				{"SNMP community", *conf.SNMPCommunity},
				{"SNMP v3 user", *conf.SNMPUser},
				{"SNMP v3 passphrase", *conf.SNMPAuthPass},
				{"SNMP v3 privacy", *conf.SNMPAuthPriv},
				{"SNMP listen", *conf.SNMPListen},
			},
		}
		tab.WriteTable(os.Stdout, &cfg)
	}
}
