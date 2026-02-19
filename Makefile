.PHONY: build test install clean

build:
	go build -o gr .

test:
	go test ./...

install:
	go install .

clean:
	rm -f gr
