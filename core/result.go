package core

type EvalResult struct {
	Error    string `json:"error"`
	Data     string `json:"data"`
	Datatype string `json:"datatype"`
}

type InspectResult struct {
	Error    string `json:"error"`
	Markdown string `json:"data"`
	Datatype string `json:"datatype"`
}
