package core

import (
	"fmt"
	"log"

	"github.com/dop251/goja"
)

func initBuiltins(vm *goja.Runtime) {
	vm.Set("log", func(arg ...any) {
		log.Println(arg...)
	})
	vm.Set("include", includeScript)
}

func includeScript(path string) {
	fmt.Println("including", path)
}
