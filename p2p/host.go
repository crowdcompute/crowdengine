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

	"github.com/crowdcompute/crowdengine/common"
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

type Host struct {
	P2PHost  host.Host
	dht      *dht.IpfsDHT
	IP       string
	FullAddr string

	*SwarmProtocol
	*TaskProtocol
	*DiscoveryProtocol
	*UploadImageProtocol
	*InspectContainerProtocol
	*ListImagesProtocol
}

// NewHost creates a new Host
func NewHost(port int, IP string, bootnodes []string) *Host {
	host := &Host{IP: IP}
	host.makeRandomHost(port, IP)

	if len(bootnodes) > 0 {
		host.ConnectWithNodes(bootnodes)
	}
	log.Print("Here is my p2p ID: ")
	host.FullAddr = fmt.Sprintf("/ip4/%s/tcp/%d/ipfs/%s", IP, port, host.P2PHost.ID().Pretty())
	log.Println(host.FullAddr)
	host.registerProtocols()
	return host
}

// Registering all Protocols
func (h *Host) registerProtocols() {
	// TODO: PATH has to be in a config
	h.SwarmProtocol = NewSwarmProtocol(h.P2PHost, h.IP)
	h.DiscoveryProtocol = NewDiscoveryProtocol(h.P2PHost, h.dht)
	h.TaskProtocol = NewTaskProtocol(h.P2PHost)
	// Registering the Observer that wants to get notified when the task is done.
	h.TaskProtocol.Register(h.DiscoveryProtocol)
	h.UploadImageProtocol = NewUploadImageProtocol(h.P2PHost)
	h.InspectContainerProtocol = NewInspectContainerProtocol(h.P2PHost)
	h.ListImagesProtocol = NewListImagesProtocol(h.P2PHost)
}

// makeRandomHost creates a libp2p host with a randomly generated identity.
func (h *Host) makeRandomHost(port int, IP string) {
	// Ignoring most errors for brevity
	// priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	// priv, _, _ := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	priv, _, err := crypto.GenerateKeyPair(crypto.Secp256k1, 256)
	common.CheckErr(err, "[makeRandomHost] Can't GenerateKeyPair ")

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
	err = h.dht.Bootstrap(ctx)
	common.CheckErr(err, "[makeRandomHost] Couldn't bootstrap the host.")
}

// ConnectWithNodes establishes a libp2p connection with the nodes
func (h *Host) ConnectWithNodes(nodes []string) {
	log.Println("Connecting to the nodes: ", nodes)
	for _, nodeAddr := range nodes {
		h.addAddrToPeerstore(nodeAddr)
	}
}

// addAddrToPeerstore parses a peer multiaddress and adds
// it to the given host's peerstore, so it knows how to
// contact it. It returns the peer ID of the remote peer.
func (h *Host) addAddrToPeerstore(addr string) peer.ID {
	// The following code extracts target's the peer ID from the
	// given multiaddress
	ipfsaddr, err := ma.NewMultiaddr(addr)
	common.CheckErr(err, "[addAddrToPeerstore] NewMultiaddr function failed.")

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	common.CheckErr(err, "[addAddrToPeerstore] ValueForProtocol function failed.")

	peerid, err := peer.IDB58Decode(pid)
	common.CheckErr(err, "[addAddrToPeerstore] IDB58Decode function failed.")
	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// We have a peer ID and a targetAddr so we add
	// it to the peerstore so LibP2P knows how to contact it
	h.P2PHost.Peerstore().AddAddr(peerid, targetAddr, ps.PermanentAddrTTL)

	return peerid
}

func (h *Host) PeerCount() int {
	return 0
}
