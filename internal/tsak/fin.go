package tsak

import (
	"github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/banner"
)

func Fin() {
	log.Debug("[ tsak2 ] tsak.Fin() is reached")
	banner.Banner("[ Zay Gezunt ]")
}
