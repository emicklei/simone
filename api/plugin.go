package api

import (
	"strings"
)

// Plugin is the required interface for plugins that expose Javascript functions to the runtime.
// Each exposed method will be available. Methods should return two values, the last being the optional error.
type Plugin interface {
	Namespace() string
	Init(ctx PluginContext) error
}

// PrintFunc is used to register a custom printer for a give typed value.
type PrintFunc func(v any, b *strings.Builder)

type PluginContext interface {
	// Set adds a value (value or function) to the Javascript VM namespace
	Set(name string, value any) error
}

type LoginHandler interface {
	Login(username, password string) error
}
