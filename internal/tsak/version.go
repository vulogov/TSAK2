package tsak

import (
	"fmt"

	"github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/banner"
	"github.com/vulogov/TSAK2/internal/conf"
)

func Version() {
	Init()
	log.Debug("[ TSAK2 ] tsak.Version() is reached")
	banner.Banner(fmt.Sprintf("[ TSAK2 %v ]", conf.BVersion))
	banner.Table()
}
