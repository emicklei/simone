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
