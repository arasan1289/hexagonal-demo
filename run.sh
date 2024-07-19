#! /bin/bash
echo "Generating swagger docs ......................"
swag fmt
swag init -g ./cmd/http/main.go

echo "Building and running the application ......................"
go build -o ./bin/app ./cmd/http/main.go
./bin/app
