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

	"github.com/crowdcompute/crowdengine/accounts/keystore"
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
	ks              *keystore.KeyStore
}

// NewNode returns new Node instance
func NewNode(cfg *config.GlobalConfig) (*Node, error) {
	n := &Node{
		cfg:  cfg,
		quit: make(chan struct{}),
		ks:   keystore.NewKeyStore(cfg.Global.KeystoreDir),
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
func (n *Node) apis() []ccrpc.API {
	return []ccrpc.API{
		{
			Namespace:    "discovery",
			Version:      "1.0",
			Service:      ccrpc.NewDiscoveryAPI(n.host),
			Public:       true,
			AuthRequired: "",
		},
		{
			Namespace:    "imagemanager",
			Version:      "1.0",
			Service:      ccrpc.NewImageManagerAPI(n.host),
			Public:       true,
			AuthRequired: "UploadImage",
		},
		{
			Namespace:    "service",
			Version:      "1.0",
			Service:      ccrpc.NewSwarmServiceAPI(n.host),
			Public:       true,
			AuthRequired: "",
		},
		{
			Namespace:    "bootnodes",
			Version:      "1.0",
			Service:      ccrpc.NewBootnodesAPI(n.host),
			Public:       true,
			AuthRequired: "",
		},
		{
			Namespace:    "container",
			Version:      "1.0",
			Service:      ccrpc.NewContainerService(),
			Public:       true,
			AuthRequired: "",
		},
		{
			Namespace:    "image",
			Version:      "1.0",
			Service:      ccrpc.NewImageService(),
			Public:       true,
			AuthRequired: "",
		},
		{
			Namespace:    "swarm",
			Version:      "1.0",
			Service:      ccrpc.NewSwarmService(),
			Public:       true,
			AuthRequired: "",
		},
		{
			Namespace:    "accounts",
			Version:      "1.0",
			Service:      ccrpc.NewAccountsAPI(n.host, n.ks),
			Public:       true,
			AuthRequired: "LockAccount",
		},
	}
}

// StartHTTP starts a http server
func (n *Node) StartHTTP() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", ccrpc.ServeHTTP(n.apis(), n.ks))
	serveMux.HandleFunc("/upload", ccrpc.ServeFilesHTTP(n.ks, n.cfg.Global.DataDir))

	port := n.cfg.RPC.HTTP.ListenPort
	log.Println("RPC listening to the port: ", port)
	httpAddr := fmt.Sprintf("%s:%d", n.cfg.RPC.HTTP.ListenAddress, port)
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
