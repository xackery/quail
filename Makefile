NAME := quail
BUILD_VERSION ?= 2.3.0
EQPATH := ~/Documents/games/EverQuest.app/drive_c/rebuildeq/

# build quail for local OS and windows
build:
	@echo "build: building to bin/quail..."
	go build main.go
	-mv main bin/quail

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
.PHONY: test
test:
	@echo "test: running tests with 30s timeout..."
	@mkdir -p test
	@rm -rf test/*
	@go test -timeout 30s ./...

.PHONY: test-all
test-all:
	@echo "test-all: running every test, even ones flagged SINGLE_TEST timeout 120s..."
	@mkdir -p test
	@rm -rf test/*
	@SINGLE_TEST=1 go test -timeout 120s -tags ci ./...

# build all supported os's
build-all: build-darwin build-windows build-linux build-windows-addon

build-darwin:
	@echo "build-darwin: ${BUILD_VERSION}"
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -buildmode=pie -ldflags="-X main.Version=${BUILD_VERSION} -s -w" -o bin/${NAME}-darwin main.go

build-linux:
	@echo "build-linux: ${BUILD_VERSION}"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-X main.Version=${BUILD_VERSION} -s -w" -o bin/${NAME}-linux main.go

build-windows:
	@echo "build-windows: ${BUILD_VERSION}"
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -buildmode=pie -ldflags="-X main.Version=${BUILD_VERSION} -s -w" -o bin/${NAME}.exe main.go

build-windows-addon:
	@echo "build-windows-addon: ${BUILD_VERSION}"
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -buildmode=pie -ldflags="-X main.Version=${BUILD_VERSION} -X main.ShowVersion=1 -s -w" -o bin/${NAME}-addon.exe main.go

# used by xackery, build darwin copy and move to blender path
build-copy: build-darwin
	@echo "copying to quail-addons..."
	cp bin/quail-darwin "/Users/xackery/Library/Application Support/Blender/3.4/scripts/addons/quail-addon/quail-darwin"

# run pprof and dump 4 snapshots of heap
profile-heap:
	@echo "profile-heap: running pprof watcher for 2 minutes with snapshots 0 to 3..."
	@-mkdir -p bin
	curl http://localhost:6060/debug/pprof/heap > bin/heap.0.pprof
	sleep 30
	curl http://localhost:6060/debug/pprof/heap > bin/heap.1.pprof
	sleep 30
	curl http://localhost:6060/debug/pprof/heap > bin/heap.2.pprof
	sleep 30
	curl http://localhost:6060/debug/pprof/heap > bin/heap.3.pprof

# peek at a heap
profile-heap-%:
	@echo "profile-heap-$*: use top20, svg, or list *word* for pprof commands, ctrl+c when done"
	go tool pprof bin/heap.$*.pprof

# run a trace on quail
profile-trace:
	@echo "profile-trace: getting trace data, this can show memory leaks and other issues..."
	curl http://localhost:6060/debug/pprof/trace > bin/trace.out
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
	@echo "VERSION=${BUILD_VERSION}.$*" >> $$GITHUB_ENV

extverdump:
	go run scripts/extverdump/main.go ../eq > scripts/extverdump/version-rof.

explore-%:
	mkdir -p test/
	rm -rf test/_*.s3d/
	rm -rf test/_*.eqg/
	go run main.go extract ../eq/$*.s3d test/_$*.s3d
	wld-cli explore test/_$*.s3d/$*.wld

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

inspect-%:
	@echo "quail inspect ../eq/$*"
	@go run main.go inspect ../eq/$*

yaml-load-%:
	@echo "quail yaml-out ../eq/$*"
	cp ../eq/$*.eqg test/$*.eqg
	@go run main.go yaml test/$*.eqg:$*.lay test/$*.lay.yaml

yaml-save-%:
	@echo "quail yaml-in ../eq/$*"
	cp test/$*.eqg test/$*-out.eqg
	@go run main.go yaml test/$*.lay.yaml test/$*.lay
	@#go run main.go yaml test/$*.lay.yaml test/$*-out.eqg:$*.lay
