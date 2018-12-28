package rpc

import (
	"flag"
	"log"
	"net/http"

	"github.com/crowdcompute/crowdengine/fileserver"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	addrHttp = flag.String("addrHTTP", "localhost:8080", "http service address")
	addrWS   = flag.String("addr", "localhost:8088", "websocket service address")
)

// StartHTTP build a jsonrpc HTTP server
func StartHTTP() {
	server := rpc.NewServer()
	conService := new(ContainerService)
	imageService := new(ImageService)
	swarmService := new(SwarmService)
	if err := server.RegisterName("container", conService); err != nil {
		log.Fatal(err)
	}

	if err := server.RegisterName("image", imageService); err != nil {
		log.Fatal(err)
	}

	if err := server.RegisterName("swarm", swarmService); err != nil {
		log.Fatal(err)
	}

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", server.ServeHTTP)
	serveMux.HandleFunc("/upload", fileserver.ServeHTTP)

	log.Fatal(http.ListenAndServe(*addrHttp, serveMux))
}

// StartWebSocket build a jsonrpc server
func StartWebSocket() {
	// server := rpc.NewServer()
	// service := new(ImageService)
	// if err := server.RegisterName("image", service); err != nil {
	// 	log.Fatal(err)
	// }
	// serveMux := http.NewServeMux()
	// serveMux.Handle("/", server.WebsocketHandler([]string{"*"}))

	// log.Fatal(http.ListenAndServe(*addrWS, serveMux))
}
