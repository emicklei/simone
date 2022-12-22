package main

import (
	"github.com/emicklei/simone"
	"github.com/emicklei/simone/plugins/fs"
)

func main() {
	cfg := simone.Config{
		Origin:   "http://localhost:5000",
		HttpAddr: ":9119",
		Plugins:  []simone.Plugin{new(Demo), new(fs.Plugin)},
	}
	simone.Start(cfg)
}
