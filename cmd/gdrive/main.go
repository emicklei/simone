package main

import "github.com/emicklei/simone"

// start flutter with options: --web-hostname=localhost --web-port=5000
func main() {
	cfg := simone.Config{
		Origin:   "http://localhost:5000",
		HttpAddr: ":8888",
	}
	gdriver := new(GDrive)
	cfg.Initializers = append(cfg.Initializers, gdriver)
	simone.Start(cfg)
}
