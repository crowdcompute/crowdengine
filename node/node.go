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

package node

import (
	"errors"
	"flag"
	"net/http"
	"sync"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/fileserver"

	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/p2p"
	ccrpc "github.com/crowdcompute/crowdengine/rpc"
	"github.com/urfave/cli"

	"github.com/ethereum/go-ethereum/rpc"
)

var (
	errNodeStarted      = errors.New("node: already started")
	errImageStoreExists = errors.New("Unable to create a new Image Store")
	httpAddr            = flag.String("httpAddr", "localhost:8080", "http service address")
	addrWS              = flag.String("addrWS", "localhost:8081", "web socket service address")
	httpFileServerAddr  = flag.String("httpFileServerAddr", "localhost:8082", "http file server address")
)

// Node represents a node
type Node struct {
	rpcAPIs         []rpc.API // list of APIs
	startOnce       sync.Once
	quit            chan struct{} // Channel used for graceful exit
	store, imgTable database.Database
	host            *p2p.Host
	cfg             *config.GlobalConfig
}

// NewNode returns new Node instance
func NewNode(cfg *config.GlobalConfig) (*Node, error) {
	n := &Node{
		cfg:  cfg,
		quit: make(chan struct{}),
	}
	n.host = p2p.NewHost(cfg.P2P.ListenPort, cfg.P2P.ListenAddress, cfg.P2P.Bootstraper.Nodes)
	return n, nil
}

// Start starts a node instance & listens to RPC calls if the flag is set
func (n *Node) Start(ctx *cli.Context) error {
	err := errNodeStarted
	n.startOnce.Do(func() {
		// TODO: Only if worker node run these two
		go n.host.DeleteDiscoveryMsgs(n.quit)
		go PruneImages(n.quit)
		err = nil // clear error above, only once.
	})

	if n.cfg.RPC.Enabled {
		if n.cfg.RPC.HTTP.Enabled {
			go n.StartHTTP()
			go n.StartFileServer()
		}
		if n.cfg.RPC.Websocket.Enabled {
			n.StartWebSocket()
		}
	}

	select {}

	return err
}

// Stop is closing down everything that the node started
func (n *Node) Stop() error {
	n.store.Close()
	close(n.quit)
	log.Println("Node stopped")
	return nil
}

// apis returns the collection of RPC descriptors this node offers.
func (n *Node) apis() []rpc.API {
	return []rpc.API{
		{
			Namespace: "discovery",
			Version:   "1.0",
			Service:   ccrpc.NewDiscoveryAPI(n.host),
			Public:    true,
		},
		{
			Namespace: "imagemanager",
			Version:   "1.0",
			Service:   ccrpc.NewImageManagerAPI(n.host),
			Public:    true,
		},
		{
			Namespace: "service",
			Version:   "1.0",
			Service:   ccrpc.NewServiceAPI(n.host),
			Public:    true,
		},
		{
			Namespace: "bootnodes",
			Version:   "1.0",
			Service:   ccrpc.NewBootnodesAPI(n.host),
			Public:    true,
		},
	}
}

// StartHTTP starts a http server
func (n *Node) StartHTTP() {
	server := rpc.NewServer()
	for _, api := range n.apis() {
		err := server.RegisterName(api.Namespace, api.Service)
		common.CheckErr(err, "[StartHTTP] Ethereum RPC could not register name.")
	}

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", server.ServeHTTP)

	log.Fatal(http.ListenAndServe(*httpAddr, serveMux))
}

// StartWebSocket starts a websocket server
func (n *Node) StartWebSocket() {
	server := rpc.NewServer()
	for _, api := range n.apis() {
		err := server.RegisterName(api.Namespace, api.Service)
		common.CheckErr(err, "[StartHTTP] Ethereum RPC could not register name.")
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/", server.WebsocketHandler([]string{"*"}))

	log.Fatal(http.ListenAndServe(*addrWS, serveMux))
}

// StartFileServer starts a file server
func (n *Node) StartFileServer() {
	serveMux := http.NewServeMux()
	log.Printf("Starting upload file server...\n")
	serveMux.HandleFunc("/upload", fileserver.ServeHTTP)
	log.Fatal(http.ListenAndServe(*httpFileServerAddr, serveMux))
}
