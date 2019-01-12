// Copyright 2018 The crowdcompute:crowdengine Authors
// This file is part of the crowdcompute:crowdengine library.
//
// The crowdcompute:crowdengine library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The crowdcompute:crowdengine library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the crowdcompute:crowdengine library. If not, see <http://www.gnu.org/licenses/>.

package rpc

import (
	"flag"
	"net/http"

	"github.com/crowdcompute/crowdengine/log"

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
