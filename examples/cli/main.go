package main

import (
	"fmt"

	"github.com/emicklei/simone"
	"github.com/emicklei/simone/api"
	"github.com/emicklei/simone/plugins/fs"
)

func main() {
	cfg := api.Config{
		Plugins: []api.Plugin{new(fs.Plugin)},
	}
	if err := simone.Run(cfg); err != nil {
		fmt.Println("err:", err)
	}
}
