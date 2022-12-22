package simone

import (
	"io"
	"net/http"
	"net/url"
)

type ActionParams struct {
	Debug  bool
	Line   string //zero-based
	Action string
	File   string
	Source string
}

func (p ActionParams) Inject(base *url.URL) {
	vals := base.Query()
	if p.Debug {
		vals.Add("debug", "true")
	}
	vals.Add("line", p.Line)
	vals.Add("action", p.Action)
	vals.Add("file", p.File)
	base.RawQuery = vals.Encode()
}

func NewActionParams(req *http.Request) ActionParams {
	base := req.URL
	body, err := io.ReadAll(req.Body)
	if err == nil {
		defer req.Body.Close()
	}
	return ActionParams{
		Debug:  base.Query().Get("debug") == "true",
		Line:   base.Query().Get("line"),
		Action: base.Query().Get("action"),
		File:   base.Query().Get("file"),
		Source: string(body),
	}
}
