package api

type EvalResult struct {
	Error    string `json:"error"`
	Datatype string `json:"datatype"`
	RawData  any    `json:"-"`
}

type InspectResult struct {
	Error string `json:"error"`
	// TODO
	IsArray     bool           `json:"isarray"`
	Scalar      any            `json:"scalar"`
	Object      map[string]any `json:"object"`
	ObjectArray []any          `json:"objectarray"`

	Datatype string `json:"datatype"`
}

type HoverResult struct {
	Error    string `json:"error"`
	Markdown string `json:"markdown"`
	Datatype string `json:"datatype"`
}
