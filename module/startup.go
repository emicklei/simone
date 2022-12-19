package module

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func Start(cfg Config) {
	cc := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.Origin},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := NewActionHandler(cfg)
	mux := http.NewServeMux()
	mux.Handle("/v1/statements", handler)
	log.Println("Serving on localhost:" + cfg.HttpAddr)
	panic(http.ListenAndServe(cfg.HttpAddr, cc.Handler(mux)))
}
