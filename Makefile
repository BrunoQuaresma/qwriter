all: build

.PHONY: build
build: $(wildcard **/*.go)
	bash ./scripts/build.sh