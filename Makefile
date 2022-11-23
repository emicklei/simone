.PHONY: pb

pb:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	--dart_out=grpc:/Users/emicklei/AndroidStudioProjects/flutter_app_mdi/lib \
    api/inspection.proto

run:
	cd cmd/kiyaspace && go run *.go