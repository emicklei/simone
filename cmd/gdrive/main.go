package main

import "github.com/emicklei/simone/module"

// start flutter with options: --web-hostname=localhost --web-port=5000
func main() {
	cfg := module.Config{
		Origin:  "http://localhost:5000",
		Apiaddr: ":9090",
	}
	gdriver := new(GDrive)
	cfg.Initializers = append(cfg.Initializers, gdriver.Start)
	module.Start(cfg)
}
