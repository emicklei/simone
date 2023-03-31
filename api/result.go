package api

type EvalResult struct {
	Error    string `json:"error"`
	Datatype string `json:"datatype"`
	RawData  any    `json:"-"`
}

type InspectResult struct {
	Error    string                `json:"error"`
	Object   map[string]any        `json:"object"`
	Paths    map[string]AccessPath `json:"paths"`
	Datatype string                `json:"datatype"`
}

type HoverResult struct {
	Error    string `json:"error"`
	Markdown string `json:"markdown"`
	Datatype string `json:"datatype"`
}

type AccessPath struct {
	Expression string
}

// NoOutputValue is a value to return from a function to prevent output.
var NoOutputValue = struct{}{}
var NoOutputValueString = "__NoOutputValueString"
