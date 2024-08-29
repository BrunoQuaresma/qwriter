all: build/qwriter

build/qwriter: $(wildcard **/*.go)
	go build -o ./build/qwriter ./cmd