build:
	packr
	go build -o bin/justblog main.go
	packr clean

clean:
	rm -rf ./bin

.PHONY: build clean