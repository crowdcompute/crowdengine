DIR := ${CURDIR}
NODENAME := p2p-node
BINARY_NAME		=	gocc
VERSION         :=	$(shell cat ./VERSION)

run: build execute

test:
	go test ./... -v

build:
	go build -ldflags "-X main.Version=$(VERSION)" -o "./build/bin/$(BINARY_NAME)" "./cmd/gocc/"
	cp "./cmd/gocc/config/config.development.toml" "./build/bin/"

deploy: build provision_script

execute:
	./build/bin/${BINARY_NAME}

run-logserver:
	sudo sysctl -w vm.max_map_count=262144
	docker run -p 5601:5601 -p 9200:9200 -p 5044:5044 -d gocc-logstack
	echo "Open: http://localhost:5601"

build-logserver:
	docker build --tag gocc-logstack ./build/gocc-logstack/

provision_script:
	./build/deploy/provisioner.sh $(NODENAME)1 $(DIR)/build/bin/$(BINARY_NAME)
	#./build/deploy/provisioner.sh $(NODENAME)2 $(DIR)/build/bin/$(BINARY_NAME)

.PHONY: test build run provision_script deploy