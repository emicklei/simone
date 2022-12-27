package api

import "github.com/dop251/goja"

type Config struct {
	// Required for CORS
	Origin string
	// Optional when running as a service
	HttpAddr string
	// To configure custom plugins
	Plugins []Plugin
	// If set then this function is called on Start after all plugins are started
	Setup RuntimeSetupFunc
}

// RuntimeSetupFunc is the signature of the Setup callback.
type RuntimeSetupFunc func(vm *goja.Runtime) error
