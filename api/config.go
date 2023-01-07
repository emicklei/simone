package api

import (
	"fmt"
	"strings"

	"github.com/dop251/goja"
)

var Debug = false

type Config struct {
	// Optional for initialization
	Script string
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

// HostPort returns host:port
func (c Config) HostPort() string {
	if strings.HasPrefix(c.HttpAddr, ":") {
		return fmt.Sprintf("localhost%s", c.HttpAddr)
	}
	return c.HttpAddr
}
