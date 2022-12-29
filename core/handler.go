package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/dop251/goja"
	"github.com/emicklei/simone/api"
)

type ActionHandler struct {
	vm     *goja.Runtime
	space  *ObjectSpace
	config api.Config
}

func NewActionHandler(cfg api.Config) *ActionHandler {
	vm := goja.New()
	//vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	space := NewObjectSpace()

	// install builtins before plugin initialization
	InitBuiltins(vm)

	// init all plugins
	for _, each := range cfg.Plugins {
		ns := each.Namespace()
		log.Println("init plugin", ns)
		vm.Set(ns, each)
		if err := each.Init(vm); err != nil {
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
	vm.Set("_toggledebug", h.ToggleDebug)
	vm.Set("_showhelp", h.ShowHelp)
	if cfg.Setup != nil {
		log.Println("custom setting up Javascript virtual machine")
		if err := cfg.Setup(vm); err != nil {
			log.Fatal(err)
		}
	}
	return h
}

func (h *ActionHandler) ShowHelp() {
	log.Println("no help yet")
}

func (h *ActionHandler) ToggleDebug() {
	if api.Debug {
		api.Debug = false
		log.Println("verbose log disabled")
	} else {
		api.Debug = true
		log.Println("verbose log enabled")
	}
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
		var plugin api.Plugin
		pluginType := reflect.TypeOf(&plugin).Elem()
		if vt.Implements(pluginType) {
			continue
		}
		filtered = append(filtered, each)
	}
	return
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
		result, err := h.vm.RunString(ap.Source)
		if err != nil {
			// syntax error
			if api.Debug {
				log.Println("RunString failed:", err.Error())
			}
			res.Error = err.Error()
		} else {
			val := result.Export()
			if err, ok := val.(error); ok {
				// evaluation error
				res.Error = err.Error()
			} else {
				// no error
				res.Data = Print(val)
				res.Datatype = fmt.Sprintf("%T", val)
			}
		}
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

func (h *ActionHandler) Run(ap ActionParams) (string, error) {
	result, err := h.vm.RunString(ap.Source)
	if err != nil {
		return "", err
	}
	return Print(result.Export()), nil
}

// https://stackoverflow.com/questions/67749752/how-to-apply-styling-and-html-tags-on-hover-message-with-vscode-api
func (h *ActionHandler) MarkdownInspectionOf(token string) (string, string, error) {
	val := h.vm.Get(token)
	if val == nil {
		return "", "", nil
	}
	gv := val.Export()
	return Print(gv), fmt.Sprintf("%T", gv), nil
}
