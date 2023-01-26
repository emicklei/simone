package simone

import (
	"flag"

	"github.com/emicklei/simone/api"
	"github.com/emicklei/simone/core"
)

var (
	// start options
	oDebug         = flag.Bool("v", false, "verbose logging")
	oStartupScript = flag.String("s", "", "run script from filename on startup")
	oRunScript     = flag.String("i", "", "run the script from filename as input")
	oClient        = flag.Bool("c", false, "start a client REPL")
	oHelp          = flag.Bool("h", false, "show help")
	oWindowUI      = flag.Bool("w", false, "open Simone Workspace")
)

// Start runs the application is one of the modes:
// - local evaluate a script from file
// - http service + REPL
// - http client + REPL
func Start(cfg api.Config) {
	flag.Parse()
	if *oHelp {
		flag.PrintDefaults()
		return
	}
	if *oDebug {
		api.Debug = true
	}
	cfg.StartupScript = *oStartupScript
	cfg.RunScript = *oRunScript
	cfg.RemoteClient = *oClient
	cfg.OpenWorkspace = *oWindowUI
	core.Start(cfg)
}

// RegisterPrinter is an alias for core impl
var RegisterPrinter = core.RegisterPrinter
