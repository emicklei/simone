package fs

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/emicklei/simone/api"
)

type Plugin struct {
	workdir string
}

func (s *Plugin) Namespace() string { return "fs" }

func (s *Plugin) Init(ctx api.PluginContext) error {
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

func (s *Plugin) WriteFile(value any, name string) error {
	var data []byte
	if s, ok := value.(string); ok {
		data = []byte(s)
	} else if d, ok := value.([]byte); ok {
		data = d
	} else {
		d, err := json.MarshalIndent(value, "", "    ")
		if err != nil {
			return err
		}
		data = d
	}
	return os.WriteFile(name, data, os.ModePerm)
}
