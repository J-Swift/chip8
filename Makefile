MAKEFLAGS += --silent

.PHONY: %

test:
	go test ./...

run:
	go run cmd/chip8/main.go $(ARGS)

build:
	@echo Building chip8
	go build -o out/chip8 cmd/chip8/main.go

clean:
	@echo Cleaning
	rm -rf out/*
