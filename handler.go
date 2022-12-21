package simone

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dop251/goja"
)

type ActionHandler struct {
	vm     *goja.Runtime
	space  *ObjectSpace
	config Config
}

func NewActionHandler(cfg Config) *ActionHandler {
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
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

func (h *ActionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ap := NewActionParams(r)
	result, err := h.vm.RunString(ap.Source)
	if err != nil {
		log.Println("RunString failed:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	log.Printf("%#v (%T)\n", result, result)
	io.WriteString(w, Print(result.Export()))
}
