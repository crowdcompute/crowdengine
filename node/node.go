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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

// authRequired is a middleware for the HTTP server.
// Authenticates a token and passes the request to the next handler
func authRequired(apis []ccrpc.API, ks *keystore.KeyStore, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if empty body
		if r.ContentLength == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		protected, err := isMethodProtected(apis, buf.Bytes())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Restore the r.Body to its original state
		r.Body = ioutil.NopCloser(buf)

		// ns is protected, place the logic which verifies the header
		if protected {
			key, err := getKeyForAccount(ks, r.Header)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), common.ContextKeyPrivateKey, key)
			log.Printf("Token valid and account {%s} unlocked. ", key.Address)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		next.ServeHTTP(w, r)
	}
}

func isMethodProtected(apis []ccrpc.API, rawJSONBody []byte) (bool, error) {
	namespace, method, err := ccrpc.FindNamespaceMethod(rawJSONBody)
	if err != nil {
		return false, err
	}
	// find which namespace
	namespaceMethodProtected := false
	for _, v := range apis {
		if v.Namespace == namespace {
			// if * then all methods are protected
			if v.AuthRequired == "*" {
				namespaceMethodProtected = true
				break
			}

			// break them and inspect them
			fncs := strings.Split(v.AuthRequired, ",")
			for _, w := range fncs {
				if common.LcFirst(strings.TrimSpace(w)) == method {
					namespaceMethodProtected = true
					break
				}
			}
			break
		}
	}
	return namespaceMethodProtected, nil
}

// Extracts the token from authorization header,
// and checks if token valid and related acount unlocked.
// And returns the key
func getKeyForAccount(ks *keystore.KeyStore, header http.Header) (*keystore.Key, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		err := fmt.Errorf("No Authorization given on header")
		log.Println(err.Error())
		return nil, err
	}
	token := strings.Split(authHeader, " ")[1]
	key, err := ks.GetKeyIfUnlockedAndValid(token)
	if err != nil {
		log.Println("Error while trying to get key for a token. Error: ", err)
		return nil, err
	}
	return key, nil
}

// UploadAuth authenticates a token and enriches the requests
// Authenticates a token and passes the request to the next handler
func uploadAuth(ks *keystore.KeyStore, uploadPath string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key, err := getKeyForAccount(ks, r.Header)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), common.ContextKeyPrivateKey, key)
		ctx = context.WithValue(ctx, common.ContextKeyUploadPath, uploadPath)
		log.Printf("Token valid and account {%s} unlocked. ", key.Address)
		next(w, r.WithContext(ctx))
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
	serveMux.Handle("/", authRequired(n.apis(), n.ks, server))
	serveMux.HandleFunc("/upload", uploadAuth(n.ks, n.cfg.Global.DataDir, ccrpc.ServeHTTP))

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
