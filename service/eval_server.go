package service

import (
	"context"

	"github.com/dop251/goja"
	"github.com/emicklei/simone/api"
)

type EvalServer struct {
	api.UnimplementedEvaluationServiceServer
	vm *goja.Runtime
}

func NewEvalServer() *EvalServer {
	return &EvalServer{
		vm: goja.New(),
	}
}

func (e *EvalServer) Eval(context.Context, *api.EvalRequest) (*api.EvalResponse, error) {
	r := new(api.EvalResponse)
	r.Result = new(api.Object)
	return r, nil
}
