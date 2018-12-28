VERSION         :=	$(shell cat ./VERSION)
BINARY_NAME		=	gocc

all: install

install:
	go install -v

build:
	go build  -o "./build/bin/$(BINARY_NAME)" -v "./cmd/cccli/main.go"

execute:
	./build/bin/${BINARY_NAME}

run: build execute

test:
	go test ./... -v

release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)
	goreleaser --rm-dist

image:
	#docker build -t {name} .

.PHONY: install test release build run