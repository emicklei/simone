package api

import (
	"reflect"
	"testing"

	"github.com/dop251/goja"
)

func createListProxy(vm *goja.Runtime) goja.Proxy {
	it := NewPagingIterator(func(input PagingState) ([]int, PagingState) {
		list := []int{1, 2, 3, 4, 5}
		index := input.NextPageIndex()
		end := index + 2
		if end > len(list) {
			end = len(list)
		}
		return list[index:end], input.WithPageIndex(end, end == len(list))
	})
	return it.ToProxy(vm)
}

func TestIterateMap(t *testing.T) {
	vm := goja.New()
	vm.Set("l", createListProxy(vm))
	_, err := vm.RunString(`r = l.map((each)=>each*each)`)
	if err != nil {
		t.Fatal(err)
	}
	gr := vm.Get("r").Export()
	if got, want := gr, []int{1, 4, 9, 16, 25}; reflect.DeepEqual(got, want) {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestIterateFilter(t *testing.T) {
	vm := goja.New()
	vm.Set("l", createListProxy(vm))
	_, err := vm.RunString(`r = l.filter((each)=>(each%2) == 0)`)
	if err != nil {
		t.Fatal(err)
	}
	gr := vm.Get("r").Export()
	if got, want := gr, []int{2, 4}; !reflect.DeepEqual(got, want) {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestIterateGetIndex(t *testing.T) {
	vm := goja.New()
	vm.Set("l", createListProxy(vm))
	_, err := vm.RunString(`r = l[3]`)
	if err != nil {
		t.Fatal(err)
	}
	gr := vm.Get("r").Export()
	if got, want := gr, 3; reflect.DeepEqual(got, want) {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestIterateLength(t *testing.T) {
	vm := goja.New()
	vm.Set("l", createListProxy(vm))
	_, err := vm.RunString(`r = l.length`)
	if err != nil {
		t.Fatal(err)
	}
	gr := vm.Get("r").Export()
	if got, want := gr, 4; reflect.DeepEqual(got, want) {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
