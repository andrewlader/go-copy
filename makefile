BINARY_NAME=go-copy

run:
	go run cmd/go-copy/main.go

install:
	go install ./cmd/go-copy

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} cmd/go-copy/main.go

run_build:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}
