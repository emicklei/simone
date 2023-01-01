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

func (h *ActionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if api.Debug {
		log.Println(r.Method, r.URL.RequestURI())
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "POST expected got "+r.Method)
		return
	}
	res := api.EvalResult{}
	w.Header().Set("content-type", "application/json")

	ap := NewActionParams(r)
	switch ap.Action {
	case "hover":
		ires := api.InspectResult{}
		md, typ, err := h.MarkdownInspectionOf(ap.Source)
		if err == nil {
			ires.Datatype = typ
			ires.Markdown = md
		} else {
			ires.Error = err.Error()
		}
		json.NewEncoder(w).Encode(ires)
		return
	case "eval", "inspect":
		if api.Debug {
			log.Println("eval", ap.Source)
		}
		res = h.runner.RunString(ap.Source)
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

// https://stackoverflow.com/questions/67749752/how-to-apply-styling-and-html-tags-on-hover-message-with-vscode-api
func (h *ActionHandler) MarkdownInspectionOf(token string) (string, string, error) {
	// val := h.vm.Get(token)
	// if val == nil {
	// 	return "", "", nil
	// }
	// gv := val.Export()
	// return Print(gv), fmt.Sprintf("%T", gv), nil
	return "TODO", "TODO", nil
}
