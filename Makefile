VERSION         :=	$(shell cat ./VERSION)
BINARY_NAME		=	gocc

DIR := ${CURDIR}

NODENAME := p2p-node

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

provision_script:
	./build/deploy/provisioner.sh $(NODENAME) $(DIR)/build/bin/$(BINARY_NAME)

deploy: build provision_script

image:
	#docker build -t {name} .

.PHONY: install test release build run provision_script deploy