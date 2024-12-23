NAME := quail
BUILD_VERSION ?= 1.4
VERSION ?= 1.4

SHELL := /bin/bash

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

lint:
	golangci-lint run

test-known:
	go test -tags ci -run TestKnown -v

##@ Build
.PHONY: build
build: ## build quail for local OS and windows
	@echo "build: building to bin/quail..."
	go build main.go
	-mv main bin/quail

.PHONY: run
run: ## run quail
	@echo "run: running..."
	go run main.go

bundle: ## bundle quail with windows icon
	@echo "if go-winres is not found, run go install github.com/tc-hib/go-winres@latest"
	@echo "bundle: setting quail icon"
	go-winres simply --icon quail.png

.PHONY: test
test: ## run tests that aren't flagged for SINGLE_TEST
	@echo "test: running tests with 30s timeout..."
	@mkdir -p test
	@rm -rf test/*
	@go test -timeout 30s ./...

.PHONY: test-all
test-all: ## test all, including SINGLE_TEST
	@echo "test-all: running every test, even ones flagged SINGLE_TEST timeout 120s..."
	@mkdir -p test
	@rm -rf test/*
	@IS_TEST_EXTENSIVE=1 SINGLE_TEST=1 go test -timeout 120s -tags ci ./...

build-all: build-darwin build-windows build-linux build-windows-addon ## build all supported os's


build-darwin: ## build darwin
	@echo "build-darwin: ${BUILD_VERSION}"
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -buildmode=pie -ldflags="-X main.Version=v${VERSION} -X main.ShowVersion=1 -s -w" -o bin/${NAME}-darwin main.go

build-linux: ## build linux
	@echo "build-linux: ${BUILD_VERSION}"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-X main.Version=v${VERSION} -X main.ShowVersion=1 -s -w" -o bin/${NAME}-linux main.go

build-windows: ## build windows
	@echo "build-windows: ${BUILD_VERSION}"
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -buildmode=pie -ldflags="-X main.Version=v${VERSION} -X main.ShowVersion=1 -s -w" -o bin/${NAME}.exe main.go

build-windows-addon: ## build windows blender addon
	@echo "build-windows-addon: ${BUILD_VERSION}"
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -buildmode=pie -ldflags="-X main.Version=v${VERSION} -X main.ShowVersion=1 -s -w" -o bin/${NAME}-addon.exe main.go

build-wasm: ## build wasm
	@echo "build-wasm: ${BUILD_VERSION}"
	@GOOS=js GOARCH=wasm go build -o bin/${NAME}.wasm main.go

build-copy: build-darwin ## used by xackery, build darwin copy and move to blender path
	@echo "copying to quail-addons..."
	cp bin/quail-darwin "/Users/xackery/Library/Application Support/Blender/3.4/scripts/addons/quail-addon/quail-darwin"

##@ Profiling

profile-heap: ## run pprof and dump 4 snapshots of heap
	@echo "profile-heap: running pprof watcher for 2 minutes with snapshots 0 to 3..."
	@-mkdir -p bin
	curl http://localhost:6060/debug/pprof/heap > bin/heap.0.pprof
	sleep 30
	curl http://localhost:6060/debug/pprof/heap > bin/heap.1.pprof
	sleep 30
	curl http://localhost:6060/debug/pprof/heap > bin/heap.2.pprof
	sleep 30
	curl http://localhost:6060/debug/pprof/heap > bin/heap.3.pprof

profile-heap-%: ## peek at a heap e.g. profile-heap-0
	@echo "profile-heap-$*: use top20, svg, or list *word* for pprof commands, ctrl+c when done"
	go tool pprof bin/heap.$*.pprof

profile-trace: ## run a trace on quail
	@echo "profile-trace: getting trace data, this can show memory leaks and other issues..."
	curl http://localhost:6060/debug/pprof/trace > bin/trace.out
	go tool trace bin/trace.out

sanitize: ## run sanitization against golang
	@echo "sanitize: checking for errors"
	rm -rf vendor/
	go vet -tags ci ./...
	test -z $(goimports -e -d . | tee /dev/stderr)
	-go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	gocyclo -over 99 .
	staticcheck -go 1.14 ./...
	go test -tags ci -covermode=atomic -coverprofile=coverage.out ./...
    coverage=`go tool cover -func coverage.out | grep total | tr -s '\t' | cut -f 3 | grep -o '[^%]*'`

# CICD triggers this
set-version-%:
	@echo "VERSION=${BUILD_VERSION}.$*" >> $$GITHUB_ENV

##@ Tools

extverdump: ## dump extensions
	source .env && go run scripts/extverdump/main.go $$EQ_PATH > scripts/extverdump/version-rof.

explore-%: ## shortcut for wld-cli to explore a file
	mkdir -p test/
	rm -rf test/_*.s3d/
	rm -rf test/_*.eqg/
	source .env && go run main.go extract $$EQ_PATH/$*.s3d test/_$*.s3d
	wld-cli explore test/_$*.s3d/$*.wld

exploreobjects-%: ## shortcut for wld-cli to explore a file
	mkdir -p test/
	rm -rf test/_*.s3d/
	rm -rf test/_*.eqg/
	source .env && go run main.go extract $$EQ_PATH/$*.s3d test/_$*.s3d
	wld-cli explore test/_$*.s3d/objects.wld


explorelights-%: ## shortcut for wld-cli to explore a file
	mkdir -p test/
	rm -rf test/_*.s3d/
	rm -rf test/_*.eqg/
	source .env && go run main.go extract $$EQ_PATH/$*.s3d test/_$*.s3d
	wld-cli explore test/_$*.s3d/lights.wld

exploretest-%: ## shortcut for wld-cli to explore a test file
	mkdir -p test/
	wld-cli explore test/$*.wld


wldcom-%: ## shortcut for WLDCOM.EXE for decoding
	mkdir -p test/
	source .env && go run main.go extract $$EQ_PATH/$*.s3d test/_$*.s3d
	-wine wldcom-patch.exe -d test/_$*.s3d/objects.wld test/_$*.s3d/objects.wldcom.wld
	-wine wldcom-patch.exe -d test/_$*.s3d/lights.wld test/_$*.s3d/lights.wldcom.wld
	wine wldcom-patch.exe -d test/_$*.s3d/$*.wld test/_$*.s3d/$*.wldcom.wld

extract-%: ## extract a file
	mkdir -p test/
	rm -rf test/_*.s3d/
	rm -rf test/_*.eqg/
	source .env && go run main.go extract $$EQ_PATH/$*.s3d test/_$*.s3d
	wld-cli extract test/_$*.s3d/$*.wld -f json test/_$*.wld.json

extractobj-%: ## extract a file
	mkdir -p test/
	rm -rf test/_*.s3d/
	rm -rf test/_*.eqg/
	source .env && go run main.go extract $$EQ_PATH/$*.s3d test/_$*.s3d
	wld-cli extract test/_$*.s3d/objects.wld -f json test/_$*.objects.wld.json

compressobj-%: ## compress a file
	mkdir -p test/
	wld-cli create test/_$*.objects.wld.json -f json test/_$*.objects.wld


extractfrag-%: ## extract a file
	mkdir -p test/
	rm -rf test/_*.s3d/
	rm -rf test/_*.eqg/
	source .env && go run main.go extract $$EQ_PATH/$*.s3d test/_$*.s3d
	wld-cli extract test/_$*.s3d/$*.wld test/_$*.s3d/

test-cover: ## test coverage %'s
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

inspect-%: ## inspect a file
	@echo "quail inspect ../eq/$*"
	@go run main.go inspect ../eq/$*

.PHONY: yaml-load-%
yaml-load-%: ## load a yaml file
	@echo "quail yaml-out ../eq/$*"
	cp ../eq/$*.eqg test/$*.eqg
	@go run main.go yaml test/$*.eqg:$*.lay test/$*.lay.yaml

yaml-save-%: ## save a yaml file
	@echo "quail yaml-in ../eq/$*"
	cp test/$*.eqg test/$*-out.eqg
	@go run main.go yaml test/$*.lay.yaml test/$*.lay
	@#go run main.go yaml test/$*.lay.yaml test/$*-out.eqg:$*.lay


biodiffwld-%:
	biodiff test/$*.src.wld test/$*.dst.wld

biodifftest:
	biodiff test/src.frag test/dst.frag

jddiff-%:
	wld-cli extract test/$*.src.wld -f json test/$*.src.json
	wld-cli extract test/$*.dst.wld -f json test/$*.dst.json
	-jd test/$*.src.json test/$*.dst.json > test/$*.diff
	code test/$*.diff

jsondiff-%:
	wld-cli extract test/$*.src.wld -f json test/$*.src.json
	wld-cli extract test/$*.dst.wld -f json test/$*.dst.json
	code -d test/$*.src.json test/$*.dst.json


jsondifffrag-%:
	wld-cli extract test/$*.src.wld -f json --fragindex $(FRAG) test/$*.src.json
	wld-cli extract test/$*.dst.wld -f json --fragindex $(FRAG) test/$*.dst.json
	code -d test/$*.src.json test/$*.dst.json