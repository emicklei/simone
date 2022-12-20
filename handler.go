package simone

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dop251/goja"
)

type ActionHandler struct {
	vm    *goja.Runtime
	space *ObjectSpace
}

func NewActionHandler(cfg Config) *ActionHandler {
	vm := goja.New()
	space := NewObjectSpace()
	for _, each := range cfg.Initializers {
		log.Println("starting plugin", each.Namespace())
		if err := each.Start(vm); err != nil {
			log.Fatal(err)
		}
	}
	return &ActionHandler{
		vm:    vm,
		space: space,
	}
}

func (h *ActionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ap := NewActionParams(r)
	result, err := h.vm.RunString(ap.Source)
	if err != nil {
		log.Println("RunString failed:", err.Error())
		fmt.Fprint(w, err.Error())
		return
	}
	log.Printf("%#v (%T)\n", result, result)
	io.WriteString(w, Print(result.Export()))
}
