package dsl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/glycerine/zygomys/zygo"
	"github.com/pieterclaerhout/go-log"

	"github.com/vulogov/TSAK2/internal/conf"
)

func InitDSL() (cfg *zygo.ZlispConfig) {
	cfg = zygo.NewZlispConfig("tsak")
	cfg.Prompt = "(tsak> "
	cfg.AfterScriptDontExit = false
	if *conf.Debug == true {
		cfg.Quiet = true
		cfg.Trace = true
		cfg.CpuProfile = filepath.Join(os.TempDir(), fmt.Sprintf("tsak_cpuprofile.%v", *conf.ID))
		cfg.MemProfile = filepath.Join(os.TempDir(), fmt.Sprintf("tsak_memprofile.%v", *conf.ID))
		log.Warnf("TSAK-script will be running in debug/verbose mode. Not recommended for production")
		log.Warnf("CPU profile stored into %v", cfg.CpuProfile)
		log.Warnf("Memory profile stored into %v", cfg.MemProfile)
	} else {
		cfg.Quiet = false
		cfg.Quiet = false
		log.Debugf("Debug mode off. OK for running TSAK-script in production")
	}
	return
}

func MakeEnvironment(cfg *zygo.ZlispConfig) (env *zygo.Zlisp) {
	env = zygo.NewZlisp()
	if cfg.CpuProfile != "" {
		f, err := os.Create(cfg.CpuProfile)
		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
		defer pprof.StopCPUProfile()
	}
	log.Debug("Running standard setup for environment")
	env.StandardSetup()
	TsakStandardSetup(cfg, env)
	TsakCustomSetup(cfg, env)
	return
}
