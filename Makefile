PROJECT=chrome-history

.PHONY: build clean

build:
	mkdir -p ./build
	go build -o build/$(PROJECT) ./cmd/

clean:
	rm -rf build/*
