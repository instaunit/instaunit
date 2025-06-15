
run: source protos
	go run ./service/cmd

source:
	protoc --go_out=service --go_opt=paths=source_relative --go-grpc_out=service --go-grpc_opt=paths=source_relative ./proto/*.proto

protos:
	protoc --descriptor_set_out=service.binpb --include_imports ./proto/*.proto

fixtures:
	protoc --descriptor_set_in ./service.binpb  --encode grpcservice.CreateUserRequest < ./data/user.prototext > ./data/user.pb
	protoc --descriptor_set_in ./service.binpb  --encode grpcservice.ListUserTypesResponse < ./data/usertypes.prototext > ./data/usertypes.pb
