NAME := quail
VERSION ?= 1.0.12
EQPATH := ~/Documents/games/EverQuest.app/drive_c/rebuildeq/

build: build-docker build-darwin
	@echo "build: running build-local..."
	@docker run \
	--rm \
	-w /src \
	-v ${PWD}:/src \
	-it quail-builder bash -c 'time make build-local'

run:
	@echo "run: running..."
	go run main.go

# CICD triggers this
.PHONY: set-variable
set-version:
	@echo "set-version: setting version to ${VERSION}"
	@echo "VERSION=${VERSION}" >> $$GITHUB_ENV

#go install github.com/tc-hib/go-winres@latest
bundle:
	@echo "bundle: setting quail icon"
	go-winres simply --icon quail.png

build-docker:
	@echo "build-docker: building docker image..."
	docker build -t quail-builder .github -f .github/build.dockerfile
build-local:
	@echo "build-local: building local..."
	@#go test ./...
	@#go test -cover ./...
	@echo "Building Linux..."
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o bin/quail-linux-${VERSION}
	@echo "Building Windows..."
	cd scripts/itdump && GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -ldflags "-X main.Version=${VERSION}" -o bin/quail-windows-${VERSION}.exe
test:
	@echo "test: running tests..."
	go test ./...
test-prep:
	@echo "test-prep: preparing tests..."
	@-#cp ${EQPATH}/obj_gears.mod mod/test/obj_gears.mod
	.PHONY: build-darwin
build-darwin:
	@echo "build-darwin: ${VERSION}"
	@GOOS=darwin GOARCH=amd64 go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o bin/${NAME}-darwin-x64 main.go
.PHONY: build-linux
build-linux:
	@echo "build-linux: ${VERSION}"
	@GOOS=linux GOARCH=amd64 go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -w" -o bin/${NAME}-linux-x64 main.go
.PHONY: build-windows
build-windows:
	@echo "build-windows: ${VERSION}"
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o bin/${NAME}-win-x64.exe main.go
	@#GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o bin/${NAME}-win.exe main.go