NAME := quail
VERSION ?= 2.1.0
EQPATH := ~/Documents/games/EverQuest.app/drive_c/rebuildeq/

# build quail for local OS and windows
build:
	@echo "build: building to bin/quail..."
	go build main.go
	-mv main bin/quail
	-mv main.exe bin/quail.exe

# run quail
run:
	@echo "run: running..."
	go run main.go


# bundle quail with windows icon
bundle:
	@echo "if go-winres is not found, run go install github.com/tc-hib/go-winres@latest"
	@echo "bundle: setting quail icon"
	go-winres simply --icon quail.png

# run tests that aren't flagged for SINGLE_TEST
test:
	@echo "test: running tests..."
	@go test ./...

# build all supported os's
build-all: build-darwin build-windows build-linux

build-darwin:
	@echo "build-darwin: ${VERSION}"
	@GOOS=darwin GOARCH=amd64 go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o bin/${NAME}-darwin main.go

build-linux:
	@echo "build-linux: ${VERSION}"
	@GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION} -s -w" -o bin/${NAME}-linux main.go

build-windows:
	@echo "build-windows: ${VERSION}"
	@GOOS=windows GOARCH=amd64 go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o bin/${NAME}-windows.exe main.go

# used by xackery, build darwin copy and move to blender path
build-copy: build-darwin
	@echo "copying to quail-addons..."
	cp bin/quail-darwin "/Users/xackery/Library/Application Support/Blender/3.4/scripts/addons/quail-addon/quail-darwin"

# run pprof and dump 3 snapshots of heap
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

# peek at a heap
profile-heap-%:
	@echo "profile-heap-$*: use top20, svg, or list *word* for pprof commands, ctrl+c when done"
	go tool pprof bin/heap.$*.pprof

# run a trace on quail
profile-trace:
	@echo "profile-trace: getting trace data, this can show memory leaks and other issues..."
	curl http://localhost:8082/debug/pprof/trace > bin/trace.out
	go tool trace bin/trace.out

# run sanitization against golang
sanitize:
	@echo "sanitize: checking for errors"
	rm -rf vendor/
	go vet -tags ci ./...
	test -z $(goimports -e -d . | tee /dev/stderr)
	-go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	gocyclo -over 99 .
	golint -set_exit_status $(go list -tags ci ./...)
	staticcheck -go 1.14 ./...
	go test -tags ci -covermode=atomic -coverprofile=coverage.out ./...
    coverage=`go tool cover -func coverage.out | grep total | tr -s '\t' | cut -f 3 | grep -o '[^%]*'`

# CICD triggers this
set-version-%:
	@echo "VERSION=${VERSION}.$*" >> $$GITHUB_ENV
