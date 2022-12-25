package simone

import (
	"strings"

	"github.com/dop251/goja"
)

type Plugin interface {
	Namespace() string
	Init(vm *goja.Runtime) error
}

type PrintFunc func(v any, b *strings.Builder)

type RuntimeSetupFunc func(vm *goja.Runtime) error
