package core

import (
	"testing"

	"github.com/dop251/goja"
)

func TestInclude(t *testing.T) {
	vm := goja.New()
	initBuiltins(vm)
	_, err := vm.RunString(`include("../examples/calc.sim") `)
	if err != nil {
		t.Fatal(err)
	}
}
