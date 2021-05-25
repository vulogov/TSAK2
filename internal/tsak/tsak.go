package tsak

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/vulogov/TSAK2/internal/conf"
)

func Main() {
	switch kingpin.MustParse(conf.App.Parse(os.Args[1:])) {
	case conf.Version.FullCommand():
		Version()
	case conf.Shell.FullCommand():
		Shell()
	}
	Fin()
}
