.PHONY: docs proto

# Command to generate documentation
docs:
	swag init --parseDependency --parseInternal -g cmd/app/main.go

# Command to generate proto files
proto:
	@if exist .\pkg\grpc\auth (rd /s /q .\pkg\grpc\auth)
	mkdir .\pkg\grpc\auth
	protoc --go_out=pkg/grpc/auth --go_opt=paths=source_relative --go-grpc_out=pkg/grpc/auth --go-grpc_opt=paths=source_relative api/grpc/auth/service.proto
	move /Y .\pkg\grpc\auth\api\grpc\auth\service.pb.go .\pkg\grpc\auth\service.pb.go
	move /Y .\pkg\grpc\auth\api\grpc\auth\service_grpc.pb.go .\pkg\grpc\auth\service_grpc.pb.go
	rd /s /q .\pkg\grpc\auth\api