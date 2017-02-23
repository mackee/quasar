.PHONY: all

all: build

build:
	go build -o _bin/quasar ./cmd/
