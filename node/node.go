package node

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/database"
	"github.com/crowdcompute/crowdengine/p2p"

	"github.com/ethereum/go-ethereum/rpc"
)

var (
	errNodeStarted      = errors.New("node: already started")
	errImageStoreExists = errors.New("Unable to create a new Image Store")
)

type Node struct {
	rpcAPIs         []rpc.API // list of APIs
	startOnce       sync.Once
	quit            chan struct{} // Channel used for graceful exit
	store, imgTable database.Database
	host            *p2p.Host
}

// New returns new Node instance.
func NewNode(port int, IP string, bootnodes []string) (*Node, error) {
	n := &Node{
		quit: make(chan struct{}),
	}
	// TODO: PATH has to be in a config
	// n.store, errImageStoreExists = database.NewDBStore(filepath.Join(common.LvlDBPath, "store.db"))
	// n.imgTable = database.NewTable(n.store, reflect.TypeOf(database.ImageLvlDB{}).Name())
	// common.CheckErr(errImageStoreExists, "[NewNode] This level db file already exists")
	n.host = p2p.NewHost(port, IP, bootnodes)
	return n, nil
}

// Start starts a node instance & listens to RPC calls if the flag is set
func (n *Node) Start(ctx context.Context, rpcFlag bool) error {
	err := errNodeStarted
	n.startOnce.Do(func() {
		go n.run(ctx)
		// TODO: Only if worker node run these two
		go n.host.DeleteDiscoveryMsgs(n.quit)
		go PruneImages(n.quit)
		err = nil // clear error above, only once.
	})

	// Start listening for file upload requests
	// Start listening for RPC calls
	if rpcFlag {
		fmt.Println("Starting RPC service")
		go n.StartHTTP() // blocks forever
		n.StartWebSocket()
	} else {
		// Block here forever and wait for IPFS stream requests
		select {}
	}
	return err
}

// Run the node
func (n *Node) run(ctx context.Context) {

}

func (n *Node) Stop() error {
	n.store.Close()
	close(n.quit)
	log.Println("Node stopped")
	return nil
}

//***************************************************************************************//
//*****************************// RPC server //*********************************//
//***************************************************************************************//

// apis returns the collection of RPC descriptors this node offers.
func (n *Node) apis() []rpc.API {
	return []rpc.API{
		{
			Namespace: "discovery",
			Version:   "1.0",
			Service:   NewDiscoveryAPI(n.host),
			Public:    true,
		},
		{
			Namespace: "imagemanager",
			Version:   "1.0",
			Service:   NewImageManagerAPI(n.host),
			Public:    true,
		},
		{
			Namespace: "service",
			Version:   "1.0",
			Service:   NewServiceAPI(n.host),
			Public:    true,
		},
		{
			Namespace: "bootnodes",
			Version:   "1.0",
			Service:   NewBootnodesAPI(n.host),
			Public:    true,
		},
	}
}

var (
	httpAddr = flag.String("httpAddr", "localhost:8080", "http service address")
	addrWS   = flag.String("addrWS", "localhost:8081", "web socket service address")
)

// StartHTTP starts serving HTTP requests
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

// StartWebSocket build a jsonrpc server
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
