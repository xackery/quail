NAME := quail

SHELL := /bin/bash

contract:
	@-rm test/*.wce
	@-rm test/*.md
	@-rm test/*.ts
	@-rm test/*.py

	@go test -run ^TestWceGenMarkdown github.com/xackery/quail/wce/def
	@mv test/latest.md ../eqemu-docs-v2/docs/client/wcemu/latest.md

	@go test -run ^TestWceGenTypescript github.com/xackery/quail/wce/def
	@mv test/*.ts ../wce-vscode/src/definition

	@go test -run ^TestWceGenPython github.com/xackery/quail/wce/def
	@mv test/*.py ../quail-addon/wce

	@mkdir -p ../wce-vscode/test/read
	@-rm -rf ../wce-vscode/test/read/*
	@cp test/*.wce ../wce-vscode/test/read

	@mkdir -p ../quail-addon/test/read
	@-rm -rf ../quail-addon/test/read/*
	@cp test/*.wce ../quail-addon/test/read

	@-rm test/*.wce

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
	@echo "test-all: Running extensive tests"
	@mkdir -p test
	@rm -rf test/*
	@source .env && EQ_PATH=$$EQ_PATH SINGLE_TEST=1 go test -timeout 120s -tags test_all ./...

build-all: build-darwin build-windows build-linux build-windows-addon build-wasm ## build all supported os's

clean: test-clear
clear: test-clear

test-clear:
	@echo "test-clear: clearing test files"
	rm -rf test/*
build-darwin: ## build darwin
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -buildmode=pie -ldflags="-X main.ShowVersion=1 -s -w" -o bin/${NAME}-darwin

build-linux: ## build linux
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-X main.ShowVersion=1 -s -w" -o bin/${NAME}-linux

build-windows: ## build windows
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -buildmode=pie -ldflags="-X main.ShowVersion=1 -s -w" -o bin/${NAME}.exe

build-windows-addon: ## build windows blender addon
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -buildmode=pie -ldflags="-X main.ShowVersion=1 -s -w" -o bin/${NAME}-addon.exe

build-wasm: ## build wasm
	@GOOS=js GOARCH=wasm go build -o wasm/quail.wasm main_wasm.go

build-blender: build-linux ## used by xackery, build darwin copy and move to blender path
	@echo "copying to quail-addons..."
	cp bin/quail-linux ~/.config/blender/4.2/scripts/addons/quail-addon

build-dll: ## build dll
	@echo "building dll..."
	CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows go build -o bin/quail.dll -buildmode=c-shared ./dll/quail_dll.go
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
	go vet ./...
	test -z $(goimports -e -d . | tee /dev/stderr)
	#go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	#gocyclo -over 120 .
	# staticcheck -go 1.24 ./...
	go test -tags ci -covermode=atomic -coverprofile=coverage.out ./...
    coverage=`go tool cover -func coverage.out | grep total | tr -s '\t' | cut -f 3 | grep -o '[^%]*'`


##@ Tools

extverdump: ## dump extensions
	source .env && go run scripts/extverdump/main.go $$EQ_PATH > scripts/extverdump/version-rof.

explore-%: ## shortcut for wld-cli to explore a file
	mkdir -p test/
	rm -rf test/_*.s3d/
	rm -rf test/_*.eqg/
	source .env && go run main.go unzip $$EQ_PATH/$*.s3d test/_$*.s3d
	wld-cli explore test/_$*.s3d/$*.wld

exploreobjects-%: ## shortcut for wld-cli to explore a file
	mkdir -p test/
	rm -rf test/_*.s3d/
	rm -rf test/_*.eqg/
	source .env && go run main.go unzip $$EQ_PATH/$*.s3d test/_$*.s3d
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

cmptreetest-%:
	go run main.go tree test/$*.src.wld test/$*.dst.wld

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
	biodiff test/src.bin test/dst.bin

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

flagfinder:
	source .env && EQ_PATH=$$EQ_PATH SCRIPT_TEST=1 go test -v -run ^TestFragFlags$$ github.com/xackery/quail/wce/wld_test

test-extensive:
	EQ_PATH=/src/eq/rof2 SINGLE_TEST=1 go test -v -timeout 0 -run ^TestDoubleConvertQuailDir$$ github.com/xackery/quail/cmd