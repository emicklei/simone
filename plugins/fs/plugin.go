package fs

import (
	"os"

	"github.com/dop251/goja"
)

type Plugin struct{}

func (s *Plugin) Namespace() string { return "fs" }

func (s *Plugin) Start(vm *goja.Runtime) error {
	vm.Set(s.Namespace(), &Plugin{})
	return nil
}

func (s *Plugin) Dir() (list []string) {
	files, err := os.ReadDir(".")
	if err != nil {
		return
	}
	for _, each := range files {
		list = append(list, each.Name())
	}
	return
}
