package main

import (
	"github.com/emicklei/simone"
	"github.com/emicklei/simone/plugins/fs"
)

func main() {
	cfg := simone.Config{
		HttpAddr: ":9119",
		Plugins:  []simone.Plugin{new(Demo), new(fs.Plugin)},
	}
	simone.Start(cfg)
}