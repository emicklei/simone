package simone

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/dop251/goja"
)

type ActionHandler struct {
	vm     *goja.Runtime
	space  *ObjectSpace
	config Config
}

func NewActionHandler(cfg Config) *ActionHandler {
	vm := goja.New()
	//vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	space := NewObjectSpace()
	for _, each := range cfg.Plugins {
		log.Println("starting plugin", each.Namespace())
		if err := each.Start(vm); err != nil {
			log.Fatal(err)
		}
	}
	h := &ActionHandler{
		vm:     vm,
		space:  space,
		config: cfg,
	}
	vm.Set("_plugins", h.PluginInfo)
	vm.Set("_variables", h.GlobalVariables)
	vm.Set("log", func(arg ...any) {
		log.Println(arg...)
	})
	if cfg.Setup != nil {
		log.Println("custom setting up Javascript vm")
		if err := cfg.Setup(vm); err != nil {
			log.Fatal(err)
		}
	}
	return h
}

func (h *ActionHandler) PluginInfo() (list []string) {
	for _, each := range h.config.Plugins {
		list = append(list, each.Namespace())
	}
	return
}

func (h *ActionHandler) GlobalVariables() (filtered []string) {
	for _, each := range h.vm.GlobalObject().Keys() {
		// skip internal var and funcs
		if strings.HasPrefix(each, "_") {
			continue
		}
		v := h.vm.Get(each)
		vt := v.ExportType()
		if vt.Kind() == reflect.Func {
			continue
		}
		// https://stackoverflow.com/questions/7132848/how-to-get-the-reflect-type-of-an-interface
		var plugin Plugin
		pluginType := reflect.TypeOf(&plugin).Elem()
		if vt.Implements(pluginType) {
			continue
		}
		filtered = append(filtered, each)
	}
	return
}

func (h *ActionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "POST expected got "+r.Method)
		return
	}
	w.Header().Set("content-type", "application/json")

	log.Println(r.URL.String())

	ap := NewActionParams(r)

	switch ap.Action {
	case "hover":
		// TODO use InspectResult
		log.Println("inspect", ap.Source)
		type markdownHolder struct {
			MarkdownString string
		}
		json.NewEncoder(w).Encode(markdownHolder{MarkdownString: h.MarkdownInspectionOf(ap.Source)})
	case "eval", "inspect":
		log.Println("eval", ap.Source)
		result, err := h.vm.RunString(ap.Source)
		if err != nil {
			log.Println("RunString failed:", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}
		// special printer
		// TODO use EvalResult
		val := result.Export()
		printed := fmt.Sprintf("(%T) %s", val, Print(val))
		log.Println(printed)
		io.WriteString(w, printed)
	default:
		io.WriteString(w, "unknown or empty action:"+ap.Action)
	}
}

// https://stackoverflow.com/questions/67749752/how-to-apply-styling-and-html-tags-on-hover-message-with-vscode-api
func (h *ActionHandler) MarkdownInspectionOf(token string) string {
	val := h.vm.Get(token)
	if val == nil {
		return ""
	}
	gv := val.Export()
	b := new(strings.Builder)
	fmt.Fprintf(b, "*%T*\n\n", gv)
	b.WriteString(Print(gv))
	return b.String()
}
