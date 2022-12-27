package api

import (
	"log"
	"strconv"

	"github.com/dop251/goja"
)

type PagingState struct {
	nextPageToken string
	hasMorePages  bool
	nextPageError error
}

func (p PagingState) NextPageIndex() int {
	index, _ := strconv.Atoi(p.nextPageToken)
	return index
}

func (p PagingState) NextPageToken() string { return p.nextPageToken }

func (p PagingState) WithNextPageToken(token string) PagingState {
	p.nextPageToken = token
	p.hasMorePages = token != ""
	return p
}

func (p PagingState) WithPageIndex(index int, done bool) PagingState {
	p.nextPageToken = strconv.Itoa(index)
	p.hasMorePages = !done
	return p
}

func (p PagingState) WithError(err error) PagingState {
	p.nextPageError = err
	return p
}

func (p PagingState) Err() error {
	return p.nextPageError
}

// NextPageFunc fetches the next chunk of data (a page) and return the result with information about the next chunk (page).
type NextPageFunc[T any] func(input PagingState) ([]T, PagingState)

type PageIterator[T any] struct {
	State  PagingState
	nexter NextPageFunc[T]
	cache  []T
}

// NewPagingIterator returns an iterator that uses a paging function get chunks of data.
func NewPagingIterator[T any](nexter NextPageFunc[T]) *PageIterator[T] {
	return &PageIterator[T]{
		State:  PagingState{hasMorePages: true},
		nexter: nexter,
	}
}

func (i *PageIterator[T]) ensureCached() {
	if len(i.cache) > 0 {
		return
	}
	// if cache is empty because of result then retry
	for i.HasNextPage() {
		for _, each := range i.NextPage() {
			i.cache = append(i.cache, each)
		}
	}
}

// Map implements JS map function
func (i *PageIterator[T]) Map(block func(each T) any) (list []any) {
	i.ensureCached()
	for _, each := range i.cache {
		list = append(list, block(each))
	}
	return
}

// Filter implements JS filter function
func (i *PageIterator[T]) Filter(block func(each T) bool) (list []T) {
	i.ensureCached()
	for _, each := range i.cache {
		if block(each) {
			list = append(list, each)
		}
	}
	return
}

func (i *PageIterator[T]) GetIndex(idx int) T {
	i.ensureCached()
	return i.cache[idx]
}

func (i *PageIterator[T]) Len() int {
	i.ensureCached()
	return len(i.cache)
}

func (i *PageIterator[T]) NextPage() []T {
	if i.State.hasMorePages {
		page, nextState := i.nexter(i.State)
		i.State = nextState
		return page
	}
	return []T{}
}

func (i *PageIterator[T]) HasNextPage() bool {
	return i.State.nextPageError == nil && i.State.hasMorePages
}

// ToProxy returns a JS proxy that dispatching collection functions to the iterator
func (i *PageIterator[T]) ToProxy(vm *goja.Runtime) goja.Proxy {
	obj := vm.NewObject()
	return vm.NewProxy(obj, &goja.ProxyTrapConfig{
		Get: func(target *goja.Object, property string, receiver goja.Value) (value goja.Value) {
			// https://medium.com/@mandeepkaur1/a-list-of-javascript-array-methods-145d09dd19a0
			switch property {
			case "map":
				return vm.ToValue(i.Map)
			case "filter":
				return vm.ToValue(i.Filter)
			case "length":
				return vm.ToValue(i.Len)
			default:
				log.Println("[simone.PageIterator] error: no such property:", property)
				return goja.Null()
			}
		},
		GetIdx: func(target *goja.Object, property int, receiver goja.Value) (value goja.Value) {
			return vm.ToValue(i.GetIndex(property))
		},
	})
}
