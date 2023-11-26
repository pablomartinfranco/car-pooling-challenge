# Makefile for car-pooling-challenge
# vim: set ft=make ts=8 noet
# Copyright Cabify.com
# Licence MIT

# Variables
# UNAME		:= $(shell uname -s)

.EXPORT_ALL_VARIABLES:

# this is godly
# https://news.ycombinator.com/item?id=11939200
.PHONY: help
help:	### this screen. Keep it first target to be default
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@# this is not tested, but prepared in advance for you, Mac drivers
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

# Targets
# 
.PHONY: debug
debug:	### Debug Makefile itself
	@echo $(UNAME)

.PHONY: run
run:	### Run the go project
	@go run cmd/main.go

.PHONY: race
race:	### Run the go project with race detector
	@go run -race cmd/main.go

.PHONY: verbose
verbose:	### Run the go project with verbose
	@go run -v cmd/main.go

.PHONY: clean
clean:	### Clean the go project
	@rm -rf bin

.PHONY: build
build:	### Build the go project
	@go build -o bin/pooling cmd/main.go

.PHONY: test
test:	### Run tests
	@go test -v ./...

.PHONY: all
all: clean build test

.PHONY: dockerize
dockerize: all	### Build docker image and run it
	@docker build -t pooling:latest .
	@docker run --rm -it pooling:latest

.PHONY: valgrind
valgrind:	### Run valgrind on the go project
	@go build -gcflags "-N -l" -o bin/pooling cmd/main.go
	@valgrind --leak-check=full --show-leak-kinds=all --track-origins=yes \
		--verbose --log-file=valgrind-out.txt ./bin/pooling

.PHONY: helgrind
helgrind:	### Run hellgrind on the go project
	@go build -gcflags "-N -l" -o bin/pooling cmd/main.go
	@valgrind --tool=helgrind --verbose --log-file=helgrind-out.txt \
		./bin/pooling

.PHONY: godoc
godoc: ### Make godoc documentation to docs/documentation.html
	@godoc -all -goroot . > docs/documentation.html
