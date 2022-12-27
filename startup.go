package simone

import (
	"errors"
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
)

// Go either starts a HTTP service or runs a script ; depending on whether a file was provided
func Go(cfg api.Config) {
	if len(os.Args) == 1 {
		log.Println("no file provided so start an HTTP server")
		Start(cfg)
	} else {
		Run(cfg)
	}
}

// Start listens for actions on a HTTP endpoint.
func Start(cfg api.Config) {
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
	handler := core.NewActionHandler(cfg)
	if len(os.Args) == 1 {
		return errors.New("missing script filename")
	}
	filename := os.Args[1]
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
