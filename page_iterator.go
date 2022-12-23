package simone

import (
	"strconv"

	"github.com/dop251/goja"
)

type PagingState struct {
	NextPageToken string
	PageSize      int
	HasMorePages  bool
	nextPageError error
}

func (p PagingState) PageIndex() int {
	index, _ := strconv.Atoi(p.NextPageToken)
	return index
}

func (p PagingState) WithNextPageToken(token string) PagingState {
	p.NextPageToken = token
	p.HasMorePages = token != ""
	return p
}

func (p PagingState) WithPageIndex(index int, done bool) PagingState {
	p.NextPageToken = strconv.Itoa(index)
	p.HasMorePages = !done
	return p
}

func (p PagingState) WithError(err error) PagingState {
	p.nextPageError = err
	return p
}

func (p PagingState) Err() error {
	return p.nextPageError
}

type NextPageFunc[T any] func(input PagingState) ([]T, PagingState)

type PageIterator[T any] struct {
	State    PagingState
	lastPage []T
	nexter   NextPageFunc[T]
}

func NewIterator[T any](pageSize int, nexter NextPageFunc[T]) *PageIterator[T] {
	return &PageIterator[T]{
		State:  PagingState{PageSize: pageSize, HasMorePages: true},
		nexter: nexter,
	}
}

// Map implements JS map function
func (i *PageIterator[T]) Map(block func(each T) any) (list []any) {
	for i.HasNextPage() {
		for _, each := range i.NextPage() {
			list = append(list, block(each))
		}
	}
	return
}

func (i *PageIterator[T]) NextPage(count ...int) []T {
	if i.State.HasMorePages {
		page, nextState := i.nexter(i.State)
		i.State = nextState
		return page
	}
	return []T{}
}

func (i *PageIterator[T]) HasNextPage() bool {
	return i.State.nextPageError == nil && i.State.HasMorePages
}

// ToProxy returns a JS proxy that dispatching collection functions to the iterator
func (i *PageIterator[T]) ToProxy(vm *goja.Runtime) goja.Proxy {
	obj := vm.NewObject()
	return vm.NewProxy(obj, &goja.ProxyTrapConfig{
		Get: func(target *goja.Object, property string, receiver goja.Value) (value goja.Value) {
			// fmt.Println("property:", property)
			return vm.ToValue(i.Map)
		},
		// GetIdx: func(target *goja.Object, property int, receiver goja.Value) (value goja.Value) {
		// 	return s.vm.ToValue(1)
		// },
	})
}
