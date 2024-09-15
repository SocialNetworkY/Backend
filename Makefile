.PHONY: docs proto

# Command to generate documentation
docs:
	swag init --parseDependency --parseInternal -g cmd/app/main.go

# Command to generate proto files
proto:
	@if exist .\pkg\gen (rd /s /q .\pkg\gen)
	@mkdir .\pkg\grpc\auth
	protoc --go_out=pkg/grpc/auth --go-grpc_out=pkg/grpc/auth --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative api/grpc/auth.proto
	@move /Y .\pkg\grpc\auth\api\grpc\auth.pb.go .\pkg\gen\auth.pb.go
	@move /Y .\pkg\grpc\auth\api\grpc\auth_grpc.pb.go .\pkg\gen\auth_grpc.pb.go
	@rd /s /q .\pkg\grpc\auth\api