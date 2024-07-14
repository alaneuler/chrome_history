PROJECT=chrome-history

.PHONY: build clean dist

build:
	mkdir -p ./build
	go build -o build/$(PROJECT) ./cmd/

clean:
	rm -rf build/*

dist: build
	cp workflow/* ./build
	cd build && zip ../chrome-history.alfredworkflow *
