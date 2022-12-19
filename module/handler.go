package module

import (
	"log"
	"net/http"

	"github.com/dop251/goja"
	"github.com/emicklei/simone"
)

type ActionHandler struct {
	vm    *goja.Runtime
	space *ObjectSpace
}

func NewActionHandler(cfg Config) *ActionHandler {
	vm := goja.New()
	space := NewObjectSpace()
	for _, each := range cfg.Initializers {
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
	ap := simone.NewActionParams(r.URL)
	log.Println(ap)
}
