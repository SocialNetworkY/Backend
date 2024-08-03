.PHONY: docs build run clean

# Command to generate documentation
docs:
	swag init --parseDependency --parseInternal -g cmd/app/main.go

# Command to build the application
build:
	go build -o bin/app cmd/app/main.go

# Command to run the application
run:
	go run cmd/app/main.go

# Command to clean up build artifacts
clean:
	rm -rf bin/
