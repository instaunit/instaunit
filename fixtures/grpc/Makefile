
build: source protos
	go build -o ./bin/service ./service/cmd

run: source protos
	go run ./service/cmd

source:
	protoc --go_out=service --go_opt=paths=source_relative --go-grpc_out=service --go-grpc_opt=paths=source_relative ./proto/*.proto

protos:
	protoc --descriptor_set_out=service.pb --include_imports ./proto/*.proto

fixtures:
	protoc --descriptor_set_in ./service.pb  --encode grpcservice.CreateUserRequest < ./data/user.prototext > ./data/user.pb
	protoc --descriptor_set_in ./service.pb  --encode grpcservice.ListUserTypesResponse < ./data/usertypes.prototext > ./data/usertypes.pb
