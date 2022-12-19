package module

import "github.com/emicklei/simone"

type Config struct {
	Origin       string
	HttpAddr     string
	Initializers []simone.Plugin
}
