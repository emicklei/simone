package api

import (
	"strings"

	"github.com/dop251/goja"
)

// Plugin is the required interface for plugins that expose Javascript functions to the runtime.
// Each exposed method will be available. Methods should return two values, the last being the optional error.
type Plugin interface {
	Namespace() string
	Init(vm *goja.Runtime) error
}

// PrintFunc is used to register a custom printer for a give typed value.
type PrintFunc func(v any, b *strings.Builder)
