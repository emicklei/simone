package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/emicklei/simone/api"
	"github.com/rs/cors"
)

func Start(cfg api.Config) {

	// run script only
	if inputFilename := cfg.RunScript; inputFilename != "" {
		r := newLocalRunner(cfg)
		res := r.Include(inputFilename)
		if res.Error != "" {
			json.NewEncoder(os.Stdout).Encode(res)
		} else {
			fmt.Println(Print(res.RawData))
		}
		return
	}

	// run client with repl
	if cfg.RemoteClient {
		r := newRemoteRunner()
		log.Println("talking to simone on localhost" + cfg.HttpAddr)
		startREPL(r)
		return
	}

	// run service with repl
	r := newLocalRunner(cfg)
	if initFilename := cfg.StartupScript; initFilename != "" {
		r.Include(initFilename)
	}
	go startHTTP(cfg, r)
	startREPL(r)
}

// startHTTP is blocking
func startHTTP(cfg api.Config, r runnable) {
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
	handler := newActionHandler(r)
	mux := http.NewServeMux()
	mux.Handle("/v1", handler)
	if cfg.HttpHandler != nil {
		log.Println("installing custom HTTP handler on \"/\"")
		mux.Handle("/", cfg.HttpHandler)
	}
	log.Println("simone is serving on localhost" + cfg.HttpAddr)
	panic(http.ListenAndServe(cfg.HttpAddr, cc.Handler(mux)))
}

// startREPL is blocking
func startREPL(r runnable) {
	cmd := newActionCommander(r)
	cmd.Loop()
}
