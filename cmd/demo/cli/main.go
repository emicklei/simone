package main

import (
	"github.com/emicklei/simone"
	"github.com/emicklei/simone/api"
	"github.com/emicklei/simone/plugins/fs"
)

func main() {
	cfg := api.Config{
		Plugins: []api.Plugin{new(fs.Plugin)},
	}
	simone.Run(cfg)
}
