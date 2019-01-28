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
	"fmt"
	"net/http"
	"sync"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/p2p"
	ccrpc "github.com/crowdcompute/crowdengine/rpc"
	"github.com/urfave/cli"

	"github.com/ethereum/go-ethereum/rpc"
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
	host, err := p2p.NewHost(cfg)
	n.host = host
	return n, err
}

// Start starts a node instance & listens to RPC calls if the flag is set
func (n *Node) Start(ctx *cli.Context) error {
	n.startOnce.Do(func() {
		// TODO: Only if worker node run these two
		go n.host.DeleteDiscoveryMsgs(n.quit)
		go PruneImages(n.quit)
	})

	if n.cfg.RPC.Enabled {
		if n.cfg.RPC.HTTP.Enabled {
			go n.StartHTTP()
		}
		if n.cfg.RPC.Websocket.Enabled {
			n.StartWebSocket()
		}
	}

	select {}
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
			Service:   ccrpc.NewSwarmServiceAPI(n.host, &n.cfg.Host.DockerSwarm),
			Public:    true,
		},
		{
			Namespace: "bootnodes",
			Version:   "1.0",
			Service:   ccrpc.NewBootnodesAPI(n.host),
			Public:    true,
		},
		{
			Namespace: "container",
			Version:   "1.0",
			Service:   new(ContainerService),
			Public:    true,
		},
		{
			Namespace: "image",
			Version:   "1.0",
			Service:   new(ImageService),
			Public:    true,
		},
		{
			Namespace: "swarm",
			Version:   "1.0",
			Service:   new(SwarmService),
			Public:    true,
		},
	}
}

// StartHTTP starts a http server
func (n *Node) StartHTTP() {
	server := rpc.NewServer()
	for _, api := range n.apis() {
		err := server.RegisterName(api.Namespace, api.Service)
		common.FatalIfErr(err, "Ethereum RPC could not register name.")
	}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", server.ServeHTTP)
	serveMux.HandleFunc("/upload", ccrpc.ServeHTTP)

	httpAddr := fmt.Sprintf("%s:%d", n.cfg.RPC.HTTP.ListenAddress, n.cfg.RPC.HTTP.ListenPort)
	log.Fatal(http.ListenAndServe(httpAddr, serveMux))
}

// StartWebSocket starts a websocket server
func (n *Node) StartWebSocket() {
	server := rpc.NewServer()
	for _, api := range n.apis() {
		err := server.RegisterName(api.Namespace, api.Service)
		common.FatalIfErr(err, "Ethereum RPC could not register name.")
	}
	serveMux := http.NewServeMux()
	serveMux.Handle("/", server.WebsocketHandler([]string{n.cfg.RPC.Websocket.CrossOriginValue}))

	addrWS := fmt.Sprintf("%s:%d", n.cfg.RPC.Websocket.ListenAddress, n.cfg.RPC.Websocket.ListenPort)
	log.Fatal(http.ListenAndServe(addrWS, serveMux))
}
