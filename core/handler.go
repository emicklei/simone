package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/emicklei/simone/api"
)

// ActionHandler handles HTTP action requests
type ActionHandler struct {
	runner Runnable
}

func NewActionHandler(r Runnable) *ActionHandler {
	return &ActionHandler{
		runner: r,
	}
}

func (h *ActionHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	ap := NewActionParams(r)
	switch ap.Action {
	case "browse":
		if api.Debug {
			log.Println("browse", ap.Source)
		}
		res := h.runner.RunString(ap.Source)
		if res.Error != "" {
			json.NewEncoder(w).Encode(res)
			return
		}
		if err := json.NewEncoder(w).Encode(res.RawData); err != nil {
			log.Println("raw data encoding failed", err)
			io.WriteString(w, err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h *ActionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	if api.Debug {
		log.Println(r.Method, r.URL.RequestURI())
	}
	if r.Method == http.MethodGet {
		h.handleGet(w, r)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "POST,GET expected got "+r.Method)
		return
	}
	res := api.EvalResult{}
	ap := NewActionParams(r)
	switch ap.Action {
	case "hover":
		if api.Debug {
			log.Println("hover", ap.Source)
		}
		// https://stackoverflow.com/questions/67749752/how-to-apply-styling-and-html-tags-on-hover-message-with-vscode-api
		res = h.runner.RunString(fmt.Sprintf("_markdowninspect(%s)", ap.Source))
		ires := api.HoverResult{
			Error:    res.Error,
			Markdown: res.RawData.(string), // not so safe
			Datatype: res.Datatype,
		}
		json.NewEncoder(w).Encode(ires)
		return
	case "eval":
		if api.Debug {
			log.Println("eval", ap.Source)
		}
		res = h.runner.RunString(ap.Source)
	case "inspect":
		if api.Debug {
			log.Println("inspect", ap.Source)
		}
		res = h.runner.RunString(ap.Source)
		ires := api.InspectResult{
			Error:    res.Error,
			Datatype: res.Datatype,
		}
		json.NewEncoder(w).Encode(ires)
		return
	default:
		if api.Debug {
			log.Println("unknown or empty action:", ap.Action)
		}
		res.Error = fmt.Sprintf("unknown or empty action:" + ap.Action)
	}
	if res.Error != "" {
		log.Println("error:", res.Error)
	}
	json.NewEncoder(w).Encode(res)
}
