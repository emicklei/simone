package main

import (
	"github.com/emicklei/simone"
	"github.com/emicklei/simone/plugins/fs"
)

func main() {
	cfg := simone.Config{
		Plugins: []simone.Plugin{new(fs.Plugin)},
	}
	simone.Run(cfg)
}
