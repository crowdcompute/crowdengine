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

package config

import (
	"github.com/urfave/cli"
)

var (
	// LogLevelFlag to show log levels
	LogLevelFlag = cli.StringFlag{
		Name:  "loglevel",
		Usage: "Logging level",
	}

	// DataDirFlag to store all the engine related data
	DataDirFlag = cli.StringFlag{
		Name:  "datadir",
		Usage: "Data sirectory to store data/metadata",
	}

	// DatabaseNameFlag used to store user data
	DatabaseNameFlag = cli.StringFlag{
		Name:  "dbname",
		Usage: "Leveldb database name",
	}

	// AvailabilityFlag defines the availability for processing
	AvailabilityFlag = cli.StringFlag{
		Name:  "availability",
		Usage: "Availability hours for processing",
	}

	// MaxContainersFlag defines max number of containers
	MaxContainersFlag = cli.IntFlag{
		Name:  "maxcontainers",
		Usage: "Maximum number of containers allowed to be run in parallel",
	}

	// CPUPerContainerFlag defines cpu numbers of each container
	CPUPerContainerFlag = cli.IntFlag{
		Name:  "containercpus",
		Usage: "Number of CPUs available to each container",
	}

	// GPUPerContainerFlag defines gpu numbers of each container
	GPUPerContainerFlag = cli.IntFlag{
		Name:  "containergpus",
		Usage: "Number of GPUs available to each container",
	}

	// MemoryPerContainerFlag defines amount of memory of a container
	MemoryPerContainerFlag = cli.IntFlag{
		Name:  "containermemory",
		Usage: "Amount of memory available to a container",
	}

	// StoragePerContainerFlag defines amount of storage of a container
	StoragePerContainerFlag = cli.IntFlag{
		Name:  "containerstorage",
		Usage: "Amount of storage available to a container",
	}

	// RPCFlag allow rpc
	RPCFlag = cli.BoolFlag{
		Name:  "rpc",
		Usage: "Enable JSON-RPC protocol",
	}

	// RPCServicesFlag list of rpc services available
	RPCServicesFlag = cli.StringFlag{
		Name:  "rpcservices",
		Usage: "List of rpc services that are allowed to be called",
	}

	// RPCWhitelistFlag allows ips to connect
	RPCWhitelistFlag = cli.StringFlag{
		Name:  "rpcwhitelist",
		Usage: "IP address whitelist",
	}

	// RPCSocketEnabledFlag indicates whether or not to listen on unix sock
	RPCSocketEnabledFlag = cli.BoolFlag{
		Name:  "socket",
		Usage: "Listen on unix domain socket",
	}

	// RPCSocketPathFlag sockpath
	RPCSocketPathFlag = cli.StringFlag{
		Name:  "socketpath",
		Usage: "Unix domain socket path",
	}

	// RPCHTTPEnabledFlag rpc over http
	RPCHTTPEnabledFlag = cli.BoolFlag{
		Name:  "http",
		Usage: "Enable RPC over HTTP",
	}

	// RPCHTTPPortFlag rpc http port
	RPCHTTPPortFlag = cli.IntFlag{
		Name:  "httpport",
		Usage: "HTTP port of RPC",
	}

	// RPCHTTPAddrFlag rpc http listen address
	RPCHTTPAddrFlag = cli.StringFlag{
		Name:  "httpaddr",
		Usage: "HTTP listen address",
	}

	// RPCHTTPCrossOriginFlag rpc http cross origin value
	RPCHTTPCrossOriginFlag = cli.StringFlag{
		Name:  "httporigin",
		Usage: "HTTP cross origin value for browser compatibility",
	}

	// RPCWSEnabledFlag rpc over ws
	RPCWSEnabledFlag = cli.BoolFlag{
		Name:  "ws",
		Usage: "Enable RPC over Websocket",
	}

	// RPCWSPortFlag rpc ws port
	RPCWSPortFlag = cli.IntFlag{
		Name:  "wsport",
		Usage: "WS port of RPC",
	}

	// RPCWSAddrFlag rpc ws listen address
	RPCWSAddrFlag = cli.StringFlag{
		Name:  "wsaddr",
		Usage: "WS listen address",
	}

	// RPCWSCrossOriginFlag rpc ws cross origin value
	RPCWSCrossOriginFlag = cli.StringFlag{
		Name:  "wsorigin",
		Usage: "WS cross origin value for browser compatibility",
	}

	// MaxPeersFlag max num of peers
	MaxPeersFlag = cli.IntFlag{
		Name:  "maxpeers",
		Usage: "Maximun number of peers allowed to connect",
	}

	// P2PListenPortFlag p2p port
	P2PListenPortFlag = cli.IntFlag{
		Name:  "port",
		Usage: "P2P port to listen",
	}

	// P2PListenAddrFlag p2p address
	P2PListenAddrFlag = cli.StringFlag{
		Name:  "addr",
		Usage: "P2P address to listen",
	}

	// P2PTimeoutFlag p2p connection timeout
	P2PTimeoutFlag = cli.IntFlag{
		Name:  "timeout",
		Usage: "P2P connection timeout between peers",
	}

	// P2PMinPeerThreasholdFlag p2p min peers threashold
	P2PMinPeerThreasholdFlag = cli.IntFlag{
		Name:  "minpeers",
		Usage: "Minimum number of peers to start periodic bootstraper",
	}

	// P2PBootstraperFlag nodes to bootstrap
	P2PBootstraperFlag = cli.StringFlag{
		Name:  "bootstrapnodes",
		Usage: "Bootstraping nodes",
	}

	// P2PPeriodicFlag bootstrap periodic timer
	P2PPeriodicFlag = cli.IntFlag{
		Name:  "bootstrapfreq",
		Usage: "Bootstraping frequency",
	}
)

// GOCCAppFlags wraps all the above together and passes to the main app
var GOCCAppFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config, c",
		Usage: "Load configuration from `FILE`",
	},
	LogLevelFlag,
	DataDirFlag,
	DatabaseNameFlag,
	AvailabilityFlag,
	MaxContainersFlag,
	CPUPerContainerFlag,
	GPUPerContainerFlag,
	MemoryPerContainerFlag,
	StoragePerContainerFlag,
	RPCFlag,
	RPCServicesFlag,
	RPCWhitelistFlag,
	RPCSocketEnabledFlag,
	RPCSocketPathFlag,
	RPCHTTPEnabledFlag,
	RPCHTTPPortFlag,
	RPCHTTPAddrFlag,
	RPCHTTPCrossOriginFlag,
	RPCWSEnabledFlag,
	RPCWSPortFlag,
	RPCWSAddrFlag,
	RPCWSCrossOriginFlag,
	MaxPeersFlag,
	P2PListenPortFlag,
	P2PListenAddrFlag,
	P2PTimeoutFlag,
	P2PMinPeerThreasholdFlag,
	P2PBootstraperFlag,
	P2PPeriodicFlag,
}
