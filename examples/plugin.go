package main

import (
	"time"

	"github.com/dop251/goja"
)

type Demo struct{}

func (d *Demo) Namespace() string { return "demo" }

func (d *Demo) Init(vm *goja.Runtime) error { return nil }

// demo.time()
func (d *Demo) Time() time.Time { return time.Now() }

type Thing struct {
	Name string
	When time.Time
}

func (d *Demo) Thing() Thing {
	return Thing{
		Name: "some",
		When: time.Now(),
	}
}
