package simone

import (
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
	oDebug          = flag.Bool("v", false, "verbose logging")
	oScript         = flag.String("s", "", "script filename")
)

func ensureFlags() {
	if flag.Parsed() {
		return
	}
	flag.Parse()
	api.Debug = *oDebug
	if api.Debug {
		log.Println("verbose logging enabled")
	}
}

// Go either starts a HTTP service or runs a script ; depending on whether a file was provided
func Go(cfg api.Config) {
	ensureFlags()
	if *oScript == "" {
		log.Println("no file provided so start an HTTP server")
		Start(cfg)
	} else {
		cfg.Script = *oScript
		if err := Run(cfg); err != nil {
			log.Fatal(err)
		}
	}
}

// Start listens for actions on a HTTP endpoint.
func Start(cfg api.Config) {
	ensureFlags()
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
	handler := core.NewActionHandler(cfg)
	mux := http.NewServeMux()
	mux.Handle("/v1", handler)
	log.Println("simone serving on localhost" + cfg.HttpAddr)
	panic(http.ListenAndServe(cfg.HttpAddr, cc.Handler(mux)))
}

// Run executes the script passed as as argument
func Run(cfg api.Config) error {
	ensureFlags()
	if *oScript != "" {
		cfg.Script = *oScript
	}
	handler := core.NewActionHandler(cfg)
	filename := cfg.Script
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	ap := core.ActionParams{
		Debug:  false,
		Action: "eval",
		File:   filename,
		Source: string(data),
	}
	output, err := handler.Run(ap)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}
