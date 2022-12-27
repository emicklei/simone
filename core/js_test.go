package core

import (
	"testing"

	"github.com/dop251/goja"
)

type Foo struct{}

func TestProxy(t *testing.T) {
	vm := goja.New()
	obj := vm.NewObject()
	p := vm.NewProxy(obj, &goja.ProxyTrapConfig{
		Get: func(target *goja.Object, property string, receiver goja.Value) (value goja.Value) {
			t.Log(property)
			m := &Map{
				list: []int{1, 2, 3, 4},
			}
			return vm.ToValue(m.Do)
		},
		GetIdx: func(target *goja.Object, property int, receiver goja.Value) (value goja.Value) {
			return vm.ToValue(1)
		},
	})
	vm.Set("p", p)
	_, err := vm.RunString(`r = p.map((each)=>each*each)`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vm.Get("r").Export())
	_, err = vm.RunString(`f = p[0]`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vm.Get("f").Export())
	_, err = vm.RunString(`l = p.length`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(vm.Get("l").Export())
}

type Map struct {
	list []int
}

func (m *Map) Do(block func(each any) any) (list []any) {
	for _, each := range m.list {
		list = append(list, block(each))
	}
	return
}

type List[T any] struct {
	target []T
}
