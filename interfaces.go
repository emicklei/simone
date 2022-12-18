package simone

import "github.com/dop251/goja"

type Plugin interface {
	Namespace() string
	Start(vm *goja.Runtime) error
}
