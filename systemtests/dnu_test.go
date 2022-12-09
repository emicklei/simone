package systemtests

import (
	"testing"

	"github.com/dop251/goja"
)

type Foo struct{}

func TestDNU(t *testing.T) {
	vm := goja.New()
	vm.Set("f", new(Foo))
	v, err := vm.RunString("f.Bar()")
	t.Log(v, err)
}
