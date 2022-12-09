package service

import (
	"context"
	"log"

	"github.com/emicklei/simone/api"
)

type InspectServer struct {
	api.UnimplementedInspectServiceServer
	space *ObjectSpace
}

func NewInspectServer(s *ObjectSpace) *InspectServer {
	return &InspectServer{
		space: s,
	}
}

func (i *InspectServer) Inspect(ctx context.Context, req *api.InspectRequest) (*api.InspectResponse, error) {
	log.Println("Inspect", req.Id)
	resp := new(api.InspectResponse)

	v := i.space.Get(req.Id)
	if v == nil {
	}

	return resp, nil
}

func (i *InspectServer) PrintString(ctx context.Context, req *api.PrintStringRequest) (*api.PrintStringResponse, error) {
	log.Println("PrintString", req.Id)
	resp := new(api.PrintStringResponse)

	v := i.space.Get(req.Id)
	if v == nil {
		resp.Content = "null"
	} else {
		resp.Content = v.String()
	}
	return resp, nil
}
