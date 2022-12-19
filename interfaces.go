package simone

import (
	"net/url"

	"github.com/dop251/goja"
)

type Plugin interface {
	Namespace() string
	Start(vm *goja.Runtime) error
}

type ActionParams struct {
	Debug  bool
	Line   string //zero-based
	Action string
	File   string
}

func (p ActionParams) Inject(base *url.URL) {
	if p.Debug {
		base.Query().Add("debug", "true")
	}
	base.Query().Add("line", p.Line)
	base.Query().Add("action", p.Action)
	base.Query().Add("file", p.File)
}

func NewActionParams(base *url.URL) ActionParams {
	return ActionParams{
		Debug:  base.Query().Get("debug") == "true",
		Line:   base.Query().Get("line"),
		Action: base.Query().Get("action"),
		File:   base.Query().Get("file"),
	}
}
