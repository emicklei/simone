package simone

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/emicklei/simone/api"
	"github.com/emicklei/simone/core"
	"github.com/rs/cors"
)

var (
	RegisterPrinter = core.RegisterPrinter
	// start options
	oDebug         = flag.Bool("v", false, "verbose logging")
	oStartupScript = flag.String("s", "", "run script from filename on startup")
	oRunScript     = flag.String("i", "", "run the script from filename as input")
	oClient        = flag.Bool("c", false, "start a client REPL")
	oHelp          = flag.Bool("h", false, "show help")
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

	// run script only
	if inputFilename := *oRunScript; inputFilename != "" {
		r := core.NewLocalRunner(cfg)
		res := r.Include(inputFilename)
		if res.Error != "" {
			json.NewEncoder(os.Stdout).Encode(res)
		} else {
			fmt.Println(core.Print(res.RawData))
		}
		return
	}

	// run client with repl
	if *oClient {
		r := core.NewRemoteRunner()
		log.Println("talking to simone on localhost" + cfg.HttpAddr)
		startREPL(r)
		return
	}

	// run service with repl
	r := core.NewLocalRunner(cfg)
	if initFilename := *oStartupScript; initFilename != "" {
		r.Include(initFilename)
	}
	go startHTTP(cfg, r)
	startREPL(r)
}

// startHTTP is blocking
func startHTTP(cfg api.Config, r core.Runnable) {
	cc := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.Origin},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	// use the address that the vscode extension in using by default
	if cfg.HttpAddr == "" {
		cfg.HttpAddr = ":9119"
	}
	handler := core.NewActionHandler(r)
	log.Println("simone is serving on localhost" + cfg.HttpAddr)
	mux := http.NewServeMux()
	mux.Handle("/v1", handler)
	panic(http.ListenAndServe(cfg.HttpAddr, cc.Handler(mux)))
}

// startREPL is blocking
func startREPL(r core.Runnable) {
	cmd := core.NewActionCommander(r)
	cmd.Loop()
}
