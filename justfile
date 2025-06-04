# Build the application
build:
    go build -o main cmd/main.go

# Run the application
run:
    go run cmd/main.go

# Run tests
test:
    go test -v ./...

# Watch for changes and rebuild
watch:
    air
