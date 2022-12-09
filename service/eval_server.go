package service

import (
	"context"
	"log"

	"github.com/dop251/goja"
	"github.com/emicklei/simone/api"
)

type VMInitializer func(vm *goja.Runtime) error

type EvalServer struct {
	api.UnimplementedEvaluationServiceServer
	vm    *goja.Runtime
	space *ObjectSpace
}

func NewEvalServer(s *ObjectSpace) *EvalServer {
	return &EvalServer{
		vm:    goja.New(),
		space: s,
	}
}

func (e *EvalServer) Initialize(extension VMInitializer) {
	extension(e.vm)
}

func (e *EvalServer) Eval(ctx context.Context, req *api.EvalRequest) (*api.EvalResponse, error) {
	log.Println("Eval", req.Source)
	result, err := e.vm.RunString(req.Source)
	if err != nil {
		return nil, err
	}
	log.Printf("%#v (%T)\n", result, result)
	id := e.space.Put(result)
	r := new(api.EvalResponse)
	r.RemoteObjectId = id
	return r, nil
}
