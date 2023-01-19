package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dop251/goja"
)

// Debug can be set on startup or by REPL command ":d"
var Debug = false

type Config struct {
	// Optional for initialization
	StartupScript string
	// Optional for command line mode
	RunScript string
	// Optional for remote client mode
	RemoteClient bool
	// Required for CORS
	Origin string
	// Optional when running as a service
	HttpAddr string
	// To configure custom plugins
	Plugins []Plugin
	// If set then this function is called on Start after all plugins are started
	Setup RuntimeSetupFunc
	// If set then this function is installed on the webserver to handle requests on "/" ( /v1 is used )
	HttpHandler http.Handler
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
