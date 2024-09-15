.PHONY: docs proto

# Command to generate documentation
docs:
	swag init --parseDependency --parseInternal -g cmd/app/main.go

# Command to generate proto files
proto:
	@if exist .\pkg\gen (rd /s /q .\pkg\gen)
	@mkdir .\pkg\gen
	protoc --go_out=pkg/gen --go-grpc_out=pkg/gen --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative api/auth.proto
	@move /Y .\pkg\gen\api\auth.pb.go .\pkg\gen\auth.pb.go
	@move /Y .\pkg\gen\api\auth_grpc.pb.go .\pkg\gen\auth_grpc.pb.go
	@rd /s /q .\pkg\gen\api