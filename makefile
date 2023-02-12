BINARY_NAME=go-copy

run:
	go run cmd/${BINARY_NAME}/${BINARY_NAME}.go

test:
	cd configs
	go run cmd/${BINARY_NAME}/${BINARY_NAME}.go -operation test -pause

install:
	go install ./cmd/${BINARY_NAME}

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} cmd/${BINARY_NAME}/${BINARY_NAME}.go

run_build:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}
