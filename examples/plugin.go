package main

import (
	"fmt"
	"time"

	"github.com/emicklei/simone/api"
)

type Demo struct{}

func (d *Demo) Namespace() string { return "demo" }

func (d *Demo) Init(ctx api.PluginContext) error { return nil }

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

func (d *Demo) Panic() {
	panic("stay calm")
}

// Implements LoginHandler to be called with `:l demo`
func (d *Demo) Login(username, password string) error {
	fmt.Println("[demo] logging with username:", username)
	return nil
}
