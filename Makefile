DIR := ${CURDIR}
NODENAME := p2p-node
BINARY_NAME		=	gocc
VERSION         :=	$(shell cat ./VERSION)

test:
	go test ./... -v

build:
	go build -ldflags "-X main.Version=$(VERSION)" -o "./build/bin/$(BINARY_NAME)" "./cmd/gocc/"
	go build -ldflags "-X main.Version=$(VERSION)" -o "./build/bin/ccpush" "./cmd/ccpush/"
	cp "./cmd/gocc/config/config.development.toml" "./build/bin/"

deploy: build provision_script

run-logserver:
	sudo sysctl -w vm.max_map_count=262144
	docker run -p 5601:5601 -p 9200:9200 -p 5044:5044 --name gocc-logstack  -d gocc-logstack 
	echo "Open: http://localhost:5601"

stop-logserver:
	docker container stop gocc-logstack
	docker container rm -f gocc-logstack

build-logserver:
	docker build --tag gocc-logstack ./build/gocc-logstack/

provision_script:
	./build/deploy/provisioner.sh $(NODENAME)1 $(DIR)/build/bin/$(BINARY_NAME)
	# ./build/deploy/provisioner.sh $(NODENAME)2 $(DIR)/build/bin/$(BINARY_NAME)

.PHONY: test build provision_script deploy