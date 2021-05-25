package tsak

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
	"github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/conf"
	"github.com/vulogov/TSAK2/internal/dsl"
	"github.com/vulogov/TSAK2/internal/signal"
)

func Shell() {
	Init()
	var fnHistory = filepath.Join(os.TempDir(), fmt.Sprintf(".tsak_history.%v", *conf.ID))
	log.Debug("[ tsak2 ] tsak.Shell() is reached")
	log.Debugf("[ tsak2 ] Shell history is stored in %v", fnHistory)

	line := liner.NewLiner()
	line.SetCompleter(func(line string) (c []string) {
		for _, n := range dsl.Completer {
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
			}
		}
		return
	})
	if f, err := os.Open(fnHistory); err == nil {
		line.ReadHistory(f)
		f.Close()
	} else {
		log.Errorf("[ tsak2 ] Error accessing to a shell history file:  %v", fnHistory)
	}
	if value, err := line.Prompt("> "); err == nil {
		log.Debug("Got: ", value)
		line.AppendHistory(value)
	} else if err == liner.ErrPromptAborted {
		log.Error("Aborted")
		signal.ExitRequest()
	} else {
		log.Error("Error reading line: ", err)
		signal.ExitRequest()
	}

	if f, err := os.Create(fnHistory); err != nil {
		log.Errorf("Error writing history file: %v", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}
