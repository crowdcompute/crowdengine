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

// Global related configuration
type Global struct {
	LogLevel     string
	DataDir      string
	KeystoreDir  string
	UploadsDir   string
	DatabaseName string
	Availability []string
}

// Host related configuration
type Host struct {
	MaxContainers       int
	CPUPerContainer     int
	GPUPerContainer     int
	MemoryPerContainer  int
	StoragePerContainer int
	DockerSwarm         DockerSwarm

	Network struct {
		IP string
	}
}

// DockerSwarm configuration.
type DockerSwarm struct {
	ListenAddress string
	ListenPort    int
}

// RPC related configuration
type RPC struct {
	Enabled         bool
	Whitelist       []string
	EnabledServices []string

	// HTTP/WS/IPC
	HTTP      HTTPWsConfig
	Websocket HTTPWsConfig
	Socket    DomainSocket
}

// P2P related configuration
type P2P struct {
	MaxPeers           int
	ListenPort         int
	ListenAddress      string
	ConnectionTimeout  int
	MinPeersThreashold int
	Bootstraper        Bootstraper
}

// DomainSocket unix/pipe socket file
type DomainSocket struct {
	Enabled bool
	Path    string
}

// Bootstraper related config
type Bootstraper struct {
	Nodes             []string
	BootstrapPeriodic int
}

// GlobalConfig wraps all the configs
type GlobalConfig struct {
	Global Global
	Host   Host
	RPC    RPC
	P2P    P2P
}

// HTTPWsConfig ...
type HTTPWsConfig struct {
	Enabled          bool
	ListenPort       int
	ListenAddress    string
	CrossOriginValue string
}
