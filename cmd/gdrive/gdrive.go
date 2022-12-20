package main

import "github.com/dop251/goja"

// https://developers.google.com/drive/api/quickstart/go
type GDrive struct{}

func (g *GDrive) Namespace() string { return "gdrive" }

func (g *GDrive) Start(vm *goja.Runtime) error { return nil }
