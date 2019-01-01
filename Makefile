DIR := ${CURDIR}
NODENAME := p2p-node
BINARY_NAME		=	gocc
VERSION         :=	$(shell cat ./VERSION)

build:
	go build -ldflags "-X main.Version=$(VERSION)" -o "./build/bin/$(BINARY_NAME)" "./cmd/gocc/main.go"

execute:
	./build/bin/${BINARY_NAME}

run: build execute

test:
	go test ./... -v

provision_script:
	./build/deploy/provisioner.sh $(NODENAME) $(DIR)/build/bin/$(BINARY_NAME)

deploy: build provision_script

image:
	#docker build -t {name} .

.PHONY: test build run provision_script deploy