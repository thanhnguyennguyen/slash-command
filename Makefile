GOBIN = $(shell pwd)/build/bin

all:
	go build -o $(GOBIN)/slashCommand ./