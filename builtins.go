package simone

import (
	"fmt"
	"log"

	"github.com/dop251/goja"
)

func InitBuiltins(vm *goja.Runtime) {
	vm.Set("log", func(arg ...any) {
		log.Println(arg...)
	})
	vm.Set("include", Include)
}

func Include(path string) {
	fmt.Println("including", path)
}
