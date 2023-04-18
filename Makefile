NAME := quail
VERSION ?= 2.0.4
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
	@#go test ./...
	EQ_PATH=/Users/xackery/Documents/games/EverQuest.app/Contents/Resources/drive_c/rebuildeq/ go test ./...
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

build-copy: build-darwin
	@echo "copying to quail-addons..."
	cp bin/quail-darwin-x64 "/Users/xackery/Library/Application Support/Blender/3.4/scripts/addons/quail-addon/quail-darwin"

profile-heap:
	@echo "profile-heap: running pprof watcher for 2 minutes with snapshots 0 to 3..."
	@-mkdir -p bin
	curl http://localhost:8082/debug/pprof/heap > bin/heap.0.pprof
	sleep 30
	curl http://localhost:8082/debug/pprof/heap > bin/heap.1.pprof
	sleep 30
	curl http://localhost:8082/debug/pprof/heap > bin/heap.2.pprof
	sleep 30
	curl http://localhost:8082/debug/pprof/heap > bin/heap.3.pprof

profile-heap-%:
	@echo "profile-heap-$*: use top20, svg, or list *word* for pprof commands, ctrl+c when done"
	go tool pprof bin/heap.$*.pprof

profile-trace:
	@echo "profile-trace: getting trace data, this can show memory leaks and other issues..."
	curl http://localhost:8082/debug/pprof/trace > bin/trace.out
	go tool trace bin/trace.out