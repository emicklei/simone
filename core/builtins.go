package core

import (
	"fmt"
	"log"

	"github.com/dop251/goja"
)

func InitBuiltins(vm *goja.Runtime) {
	log.Println("add function `log(string...)`")
	vm.Set("log", func(arg ...any) {
		log.Println(arg...)
	})
	log.Println("add function `include(string)`")
	vm.Set("include", Include)
}

func Include(path string) {
	fmt.Println("including", path)
}
