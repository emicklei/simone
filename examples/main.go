package main

import (
	"strings"
	"time"

	"github.com/emicklei/simone"
	"github.com/emicklei/simone/api"
	"github.com/emicklei/simone/plugins/fs"
)

func main() {
	cfg := api.Config{
		Plugins: []api.Plugin{new(Demo), new(fs.Plugin)},
	}
	simone.RegisterPrinter(time.Now(), printTime)

	simone.Start(cfg)
}

// custom printer for time values
func printTime(v any, b *strings.Builder) {
	t := v.(time.Time)
	b.WriteString(t.Format(time.RFC3339))
}
