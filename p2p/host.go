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

package p2p

import (
	"context"
	"fmt"

	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/log"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	ps "github.com/libp2p/go-libp2p-peerstore"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	ma "github.com/multiformats/go-multiaddr"
)

// Host represents a libp2p host
type Host struct {
	P2PHost  host.Host
	dht      *dht.IpfsDHT
	FullAddr string
	Cfg      *config.GlobalConfig

	*SwarmProtocol
	*TaskProtocol
	*DiscoveryProtocol
	*UploadImageProtocol
	*InspectContainerProtocol
	*ListImagesProtocol
}

// NewHost creates a new Host
func NewHost(cfg *config.GlobalConfig) (*Host, error) {
	nodes := cfg.P2P.Bootstraper.Nodes
	ip := cfg.P2P.ListenAddress
	port := cfg.P2P.ListenPort
	host := &Host{Cfg: cfg}
	err := host.makeRandomHost(port, ip)
	if err != nil {
		return nil, err
	}
	if len(nodes) > 0 {
		err = host.ConnectWithNodes(nodes)
	}
	log.Print("Here is my p2p ID: ")
	host.FullAddr = fmt.Sprintf("/ip4/%s/tcp/%d/ipfs/%s", ip, port, host.P2PHost.ID().Pretty())
	log.Println(host.FullAddr)
	host.registerProtocols()
	return host, err
}

// registerProtocols registers all protocols for the node
func (h *Host) registerProtocols() {
	h.SwarmProtocol = NewSwarmProtocol(h.P2PHost, h.Cfg.P2P.ListenAddress, &h.Cfg.Host.DockerSwarm)
	h.DiscoveryProtocol = NewDiscoveryProtocol(h.P2PHost, h.dht)
	h.TaskProtocol = NewTaskProtocol(h.P2PHost)
	// Registering the Observer that wants to get notified when the task is done.
	h.TaskProtocol.Register(h.DiscoveryProtocol)
	h.UploadImageProtocol = NewUploadImageProtocol(h.P2PHost)
	h.InspectContainerProtocol = NewInspectContainerProtocol(h.P2PHost)
	h.ListImagesProtocol = NewListImagesProtocol(h.P2PHost)
}

// makeRandomHost creates a libp2p host with a randomly generated identity.
func (h *Host) makeRandomHost(port int, IP string) error {
	// Ignoring most errors for brevity
	// priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	// priv, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	priv, _, err := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	if err != nil {
		return err
	}
	// listen, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", IP, port))
	listen, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	host, _ := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(listen),
		libp2p.Identity(priv),
	)

	// Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
	dstore := dsync.MutexWrap(ds.NewMapDatastore())

	// Make the DHT
	ctx := context.Background()
	h.dht = dht.NewDHT(ctx, host, dstore)

	// Make the routed host
	h.P2PHost = rhost.Wrap(host, h.dht)

	// Bootstrap the host
	return h.dht.Bootstrap(ctx)
}

// ConnectWithNodes establishes a libp2p connection with the nodes
func (h *Host) ConnectWithNodes(nodes []string) error {
	log.Println("Connecting to the nodes: ", nodes)
	for _, nodeAddr := range nodes {
		if err := h.addAddrToPeerstore(nodeAddr); err != nil {
			return err
		}
	}
	return nil
}

// addAddrToPeerstore parses a peer multiaddress and adds
// it to the given host's peerstore, so it knows how to
// contact it. It returns the peer ID of the remote peer.
func (h *Host) addAddrToPeerstore(addr string) error {
	// The following code extracts target's the peer ID from the
	// given multiaddress
	ipfsaddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		return err
	}

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		return err
	}

	peerid, err := peer.IDB58Decode(pid)
	if err != nil {
		return err
	}
	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add
	// it to the peerstore so LibP2P knows how to contact it
	h.P2PHost.Peerstore().AddAddr(peerid, targetAddr, ps.PermanentAddrTTL)

	return nil
}

// PeerCount returns the number of peers in the node's peerstore
func (h *Host) PeerCount() int {
	// Peerstore has the current node's address as well, so we don't want to count it
	return h.P2PHost.Peerstore().Peers().Len() - 1
}
