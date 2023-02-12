BINARY_NAME := go-copy
LOCAL_BINARY := ./$(BINARY_NAME)
GIT_DESCRIPTION_TXT_FILE := cmd/$(BINARY_NAME)/git-describe.txt



pre-build:
	printf "`git describe --long`\n`go version`\n`date --rfc-3339=seconds`\n" | tee cmd/go-copy/git-describe.txt

run:
	go run cmd/$(BINARY_NAME)/$(BINARY_NAME).go

test:
	cd configs
	go run cmd/$(BINARY_NAME)/$(BINARY_NAME).go -operation test -pause

install: pre-build
	go install ./cmd/$(BINARY_NAME)

build: pre-build
	GOARCH=amd64 GOOS=linux go build -o $(BINARY_NAME) cmd/$(BINARY_NAME)/$(BINARY_NAME).go

run_build:
	./$(BINARY_NAME)

build_and_run: build run_build

clean:
	go clean
	rm -f $(LOCAL_BINARY)
	rm -f $(GIT_DESCRIPTION_TXT_FILE)
