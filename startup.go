package simone

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

// Start listens for actions on a HTTP endpoint.
func Start(cfg Config) {
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
	handler := NewActionHandler(cfg)
	mux := http.NewServeMux()
	mux.Handle("/v1", handler)
	log.Println("simone serving on localhost" + cfg.HttpAddr)
	panic(http.ListenAndServe(cfg.HttpAddr, cc.Handler(mux)))
}

// Run executes the script passed as as argument
func Run(cfg Config) error {
	handler := NewActionHandler(cfg)
	if len(os.Args) == 1 {
		return errors.New("missing script filename")
	}
	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	ap := ActionParams{
		Debug:  false,
		Action: "eval",
		File:   filename,
		Source: string(data),
	}
	output, err := handler.run(ap)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}
