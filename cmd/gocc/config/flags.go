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
	"path/filepath"

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
		Value: DefaultDataDir(),
		Usage: "Data directory to store data/metadata",
	}

	// KeystoreDirFlag to store all the engine related data
	KeystoreDirFlag = cli.StringFlag{
		Name:  "keystoredir",
		Value: filepath.Join(DefaultDataDir(), "keystore"),
		Usage: "Directory for the keystore",
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

	// DockerSwarmAddrFlag defines the docker swarm's listen address
	DockerSwarmAddrFlag = cli.StringFlag{
		Name:  "swarmaddr",
		Usage: "Listen address for the docker swarm",
	}

	// DockerSwarmPortFlag defines the
	DockerSwarmPortFlag = cli.IntFlag{
		Name:  "swarmport",
		Usage: "Listen port number for the docker swarm",
	}

	// RPCFlag allow rpc
	RPCFlag = cli.BoolFlag{
		Name:  "rpc",
		Usage: "Enable JSON-RPC protocol",
	}

	// RPCServicesFlag list of rpc services available
	RPCServicesFlag = cli.StringFlag{
		Name:  "rpcservices",
		Usage: "List of rpc services allowed",
	}

	// RPCWhitelistFlag allows ips to connect
	RPCWhitelistFlag = cli.StringFlag{
		Name:  "rpcwhitelist",
		Usage: "Allow IP addresses to access the RPC servers",
	}

	// RPCSocketEnabledFlag indicates whether or not to listen on unix sock
	RPCSocketEnabledFlag = cli.BoolFlag{
		Name:  "socket",
		Usage: "Enable IPC-RPC interface",
	}

	// RPCSocketPathFlag sockpath
	RPCSocketPathFlag = cli.StringFlag{
		Name:  "socketpath",
		Usage: "Path of the socker/pipe file",
	}

	// RPCHTTPEnabledFlag rpc over http
	RPCHTTPEnabledFlag = cli.BoolFlag{
		Name:  "http",
		Usage: "Enable the HTTP-RPC server",
	}

	// RPCHTTPPortFlag rpc http port
	RPCHTTPPortFlag = cli.IntFlag{
		Name:  "httpport",
		Usage: "HTTP-RPC server listening port",
	}

	// RPCHTTPAddrFlag rpc http listen address
	RPCHTTPAddrFlag = cli.StringFlag{
		Name:  "httpaddr",
		Usage: "HTTP-RPC server listening interface",
	}

	// RPCHTTPCrossOriginFlag rpc http cross origin value
	RPCHTTPCrossOriginFlag = cli.StringFlag{
		Name:  "httporigin",
		Usage: "HTTP-RPC cross-origin value",
	}

	// RPCWSEnabledFlag rpc over ws
	RPCWSEnabledFlag = cli.BoolFlag{
		Name:  "ws",
		Usage: "Enable the WS-RPC server",
	}

	// RPCWSPortFlag rpc ws port
	RPCWSPortFlag = cli.IntFlag{
		Name:  "wsport",
		Usage: "WS-RPC server listening port",
	}

	// RPCWSAddrFlag rpc ws listen address
	RPCWSAddrFlag = cli.StringFlag{
		Name:  "wsaddr",
		Usage: "WS-RPC server listening interface",
	}

	// RPCWSCrossOriginFlag rpc ws cross origin value
	RPCWSCrossOriginFlag = cli.StringFlag{
		Name:  "wsorigin",
		Usage: "WS-RPC cross-origin value",
	}

	// MaxPeersFlag max num of peers
	MaxPeersFlag = cli.IntFlag{
		Name:  "maxpeers",
		Usage: "Maximum number of peers to connect",
	}

	// P2PListenPortFlag p2p port
	P2PListenPortFlag = cli.IntFlag{
		Name:  "port",
		Usage: "P2P listening port",
	}

	// P2PListenAddrFlag p2p address
	P2PListenAddrFlag = cli.StringFlag{
		Name:  "addr",
		Usage: "P2P listening interface",
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
	KeystoreDirFlag,
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
