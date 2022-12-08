package service

import (
	"context"
	"fmt"

	"github.com/dop251/goja"
	"github.com/emicklei/simone/api"
)

type VMInitializer func(vm *goja.Runtime) error

type EvalServer struct {
	api.UnimplementedEvaluationServiceServer
	vm *goja.Runtime
}

func NewEvalServer() *EvalServer {
	return &EvalServer{
		vm: goja.New(),
	}
}

func (e *EvalServer) Initialize(extension VMInitializer) {
	extension(e.vm)
}

func (e *EvalServer) Eval(ctx context.Context, req *api.EvalRequest) (*api.EvalResponse, error) {
	result, err := e.vm.RunString(req.Source)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	r := new(api.EvalResponse)
	r.Result = new(api.Object)
	// r.Result.Fields = append(r.Result.Fields, )
	return r, nil
}
