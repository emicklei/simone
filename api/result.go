package api

type EvalResult struct {
	Error    string `json:"error"`
	Data     string `json:"data"`
	Datatype string `json:"datatype"`
	RawData  any    `json:"-"`
}

type InspectResult struct {
	Error    string `json:"error"`
	Markdown string `json:"markdown"`
	Datatype string `json:"datatype"`
}
