BINARY_NAME=go-copy

run:
	go run cmd/main.go

install:
	go install cmd/main.go
	mv ${GOPATH}/bin/main ${GOPATH}/bin/go-copy

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} cmd/main.go

run_build:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}
