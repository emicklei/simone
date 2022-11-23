package service

import (
	"context"

	"github.com/emicklei/simone/api"
)

type InspectServer struct {
	api.UnimplementedInspectServiceServer
}

func (i *InspectServer) Inspect(context.Context, *api.InspectRequest) (*api.InspectResponse, error) {
	r := new(api.InspectResponse)
	r.Objects = map[string]*api.Object{
		"1": {
			Id: "1",
		},
	}
	return r, nil
}
