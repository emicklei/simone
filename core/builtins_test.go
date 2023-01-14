package core

import (
	"testing"

	"github.com/dop251/goja"
)

func TestInclude(t *testing.T) {
	vm := goja.New()
	vm.Set("include", includeScript)
	_, err := vm.RunString(`include("../lib/lib.sim") `)
	if err != nil {
		t.Fatal(err)
	}
}
