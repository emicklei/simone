package main

import "github.com/emicklei/simone/module"

// start flutter with options: --web-hostname=localhost --web-port=5000
func main() {
	module.Start("http://localhost:5000", ":9090")
}
