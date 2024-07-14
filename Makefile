PROJECT=chrome-history

.PHONY: build clean

build:
	mkdir -p ./build
	go build -o build/$(PROJECT) ./cmd/
	cp workflow/* ./build
	cd build && zip chrome-history.alfredworkflow *

clean:
	rm -rf build/*
