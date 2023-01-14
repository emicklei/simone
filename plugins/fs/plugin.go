package fs

import (
	"os"
	"path/filepath"

	"github.com/dop251/goja"
)

type Plugin struct {
	workdir string
}

func (s *Plugin) Namespace() string { return "fs" }

func (s *Plugin) Init(vm *goja.Runtime) error {
	s.workdir, _ = os.Getwd()
	return nil
}

// Dir is available as `fs.Dir`
func (s *Plugin) Dir() (list []string) {
	files, err := os.ReadDir(s.workdir)
	if err != nil {
		return
	}
	for _, each := range files {
		list = append(list, each.Name())
	}
	return
}

// Chdir is available as `fs.Chdir`
func (s *Plugin) Chdir(path string) string {
	if filepath.IsAbs(path) {
		s.workdir = path
	} else {
		s.workdir = filepath.Join(s.workdir, path)
	}
	if filepath.IsAbs(s.workdir) {
		return s.workdir
	} else {
		abs, _ := filepath.Abs(s.workdir)
		return abs
	}
}
