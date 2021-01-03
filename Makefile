.PHONY: clean all list-services list-nodes

.DEFAULT_GOAL := all

TARGETS=\
	list-services \
	list-nodes

cur   := $(shell pwd)

${TARGETS}:
	cd ${cur}/$@ && GO111MODULE=on go build -o ${cur}/dist/$@ main.go

all: $(TARGETS)

clean:
	rm -rf dist

lint:
	golint -set_exit_status $$(go list ./...)
	go vet ./...
