gocgi:
	mkdir -p bin
	go build -o bin/gocgi	./cmd

.PHONY: gocgi
