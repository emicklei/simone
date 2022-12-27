package main

import (
	"github.com/emicklei/simone"
	"github.com/emicklei/simone/api"
	"github.com/emicklei/simone/plugins/fs"
)

func main() {
	cfg := api.Config{
		HttpAddr: ":9119",
		Plugins:  []api.Plugin{new(Demo), new(fs.Plugin)},
	}
	simone.Start(cfg)
}
