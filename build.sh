#!/bin/bash
# Note: Made for Makefile usage
VERSION=1.0.0
set -e
apt update
apt install -y mesa-common-dev libgl1-mesa-dev libglu1-mesa-dev
cd /src
echo "Building Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$VERSION" -o bin/quail-linux-$VERSION 
echo "Building OSX..."
GOOS=darwin GOARCH=amd64 CC=o64-clang CXX=o64-clang++ go build -ldflags "-X main.Version=$VERSION" -o bin/quail-osx-$VERSION
echo "Building Windows..."
GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -ldflags "-X main.Version=$VERSION" -o bin/quail-windows-$VERSION.exe