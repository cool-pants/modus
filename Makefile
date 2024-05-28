BINARY_NAME=modus

build:
	go build -o ${BINARY_NAME} main.go

run: build
	./${BINARY_NAME}

clean:
	go clean

.PHONY: build run clean
