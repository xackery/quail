VERSION ?= 1.0.3
EQPATH := ~/Documents/games/EverQuest.app/drive_c/rebuildeq/

build: build-docker build-osx
	@docker run \
	--rm \
	-w /src \
	-v ${PWD}:/src \
	-it quail-builder bash -c 'time make build-local'
build-docker:
	docker build -t quail-builder .github -f .github/build.dockerfile
build-local:
	@#go test ./...
	@#go test -cover ./...
	@echo "Building Linux..."
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o bin/quail-linux-${VERSION} 
	@echo "Building Windows..."
	GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ go build -ldflags "-X main.Version=${VERSION}" -o bin/quail-windows-${VERSION}.exe
build-darwin:
	@echo "Building Darwin..."
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o bin/quail-darwin-${VERSION}
test-prep:
	@-#cp ${EQPATH}/obj_gears.mod mod/test/obj_gears.mod