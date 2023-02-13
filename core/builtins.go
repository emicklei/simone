package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dop251/goja"
)

func initBuiltins(vm *goja.Runtime) {
	vm.Set("log", func(arg ...any) {
		log.Println(arg...)
	})
	vm.Set("include", includer{vm: vm}.includeScript)
}

type includer struct {
	vm *goja.Runtime
}

func (i includer) includeScript(path string) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Printf("including %s -> %s\n", path, abspath)
	source, err := os.ReadFile(abspath)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = i.vm.RunString(string(source))
	if err != nil {
		log.Println(err.Error())
		return
	}
}
