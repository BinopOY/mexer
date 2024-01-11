#!/bin/bash

# Build the project
echo "Building the project..."
GOOS=linux GOARCH=amd64 go build -o ./build/mexer_amd64_linux
GOOS=windows GOARCH=amd64 go build -o ./build/mexer_amd64_windows.exe
