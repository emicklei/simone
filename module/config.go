package module

import "github.com/dop251/goja"

type VMInitializer func(vm *goja.Runtime) error

type Config struct {
	Origin       string
	Apiaddr      string
	Initializers []VMInitializer
}
