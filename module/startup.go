package module

import (
	"fmt"
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
	wrappedGrpc := grpcweb.WrapServer(grpcServer)

	// services
	api.RegisterInspectServiceServer(grpcServer, new(service.InspectServer))
	api.RegisterEvaluationServiceServer(grpcServer, service.NewEvalServer())

	// routing
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}
		// Fall back to other servers.
		mux.ServeHTTP(resp, req)
	})
	fmt.Println("gRPC on localhost", cfg.Apiaddr)
	http.ListenAndServe(cfg.Apiaddr, cc.Handler(mux))
}
