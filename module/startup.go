package module

import (
	"fmt"
	"net"
	"net/http"

	"github.com/emicklei/simone/api"
	"github.com/emicklei/simone/service"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

func Start(cfg Config) {
	cc := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.Origin},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	grpcServer := grpc.NewServer()

	// services
	eval := service.NewEvalServer()
	for _, each := range cfg.Initializers {
		eval.Initialize(each)
	}
	api.RegisterInspectServiceServer(grpcServer, new(service.InspectServer))
	api.RegisterEvaluationServiceServer(grpcServer, eval)

	// web routing
	wrappedGrpc := grpcweb.WrapServer(grpcServer)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}
		// Fall back to other servers.
		mux.ServeHTTP(resp, req)
	})
	fmt.Println("gRPC web on localhost:" + cfg.GrpcWebAddr)
	go func() {
		panic(http.ListenAndServe(cfg.GrpcWebAddr, cc.Handler(mux)))
	}()

	// grpc routing
	fmt.Println("gRPC on localhost" + cfg.GrpcAddr)
	lis, err := net.Listen("tcp", cfg.GrpcAddr)
	if err != nil {
		panic(err.Error())
	}
	grpcServer.Serve(lis)
}
