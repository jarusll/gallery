.PHONY: build
build: main.go
	go build

run: build
	./gallery
