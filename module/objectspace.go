package module

import (
	"github.com/dop251/goja"
	"github.com/google/uuid"
)

type ObjectSpace struct {
	objects map[string]goja.Value
}

func NewObjectSpace() *ObjectSpace {
	return &ObjectSpace{
		objects: map[string]goja.Value{},
	}
}

func (o *ObjectSpace) Put(v goja.Value) string {
	id := uuid.New().String()
	o.objects[id] = v
	return id
}

func (o *ObjectSpace) Get(id string) goja.Value {
	v, ok := o.objects[id]
	if !ok {
		return nil
	}
	return v
}
