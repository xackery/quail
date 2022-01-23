build:
	@docker run --rm -v ${PWD}:/src -it neilotoole/xcgo:latest bash /src/build.sh