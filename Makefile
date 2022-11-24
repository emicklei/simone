.PHONY: pb

pb:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	--dart_out=grpc:/Users/emicklei/AndroidStudioProjects/simone_fui/lib \
    api/inspection.proto \
	api/evaluation.proto

run:
	cd cmd/gdrive && go run *.go