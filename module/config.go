package module

import "github.com/emicklei/simone/service"

type Config struct {
	Origin       string
	GrpcWebAddr  string
	GrpcAddr     string
	Initializers []service.VMInitializer
}
