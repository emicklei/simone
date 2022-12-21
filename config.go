package simone

type Config struct {
	Origin   string
	HttpAddr string
	Plugins  []Plugin
	// If set then this function is called on Start after all plugins are started
	Setup RuntimeSetupFunc
}
