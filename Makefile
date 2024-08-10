.PHONY: build
build: main.go
	mkdir -p bin
	go build -o ./bin/gallery

run: build
	./bin/gallery
