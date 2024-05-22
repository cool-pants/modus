BINARY_NAME=modus

build:
	go build -o ${BINARY_NAME} main.go

run: build
	./${BINARY_NAME}

clean:
	rm ${BINARY_NAME}
	go clean

.PHONY: build run clean
