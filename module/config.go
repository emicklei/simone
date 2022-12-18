package module

import "github.com/emicklei/simone"

type Config struct {
	Origin       string
	GrpcWebAddr  string
	GrpcAddr     string
	Initializers []simone.Plugin
}
