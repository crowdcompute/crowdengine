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
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/naoina/toml"
	"github.com/urfave/cli"
)

// DefaultConfig creates a default config
func DefaultConfig() *GlobalConfig {
	return &GlobalConfig{
		Global: Global{
			LogLevel:     "TRACE",
			DataDir:      DefaultDataDir(),
			KeystoreDir:  filepath.Join(DefaultDataDir(), "keystore"),
			UploadsDir:   filepath.Join(DefaultDataDir(), "uploads"),
			DatabaseName: "gocc_db",
			Availability: []string{},
		},
		Host: Host{
			MaxContainers:       20,
			CPUPerContainer:     2,
			GPUPerContainer:     2,
			MemoryPerContainer:  1024,
			StoragePerContainer: 2048,
			DockerSwarm:         DockerSwarm{"127.0.0.1","0.0.0.0", 2377},
		},
		RPC: RPC{
			Enabled:         false,
			Whitelist:       []string{"localhost"},
			EnabledServices: []string{},
			HTTP: HTTPWsConfig{
				Enabled:          false,
				ListenPort:       8668,
				ListenAddress:    "0.0.0.0",
				CrossOriginValue: "*",
			},
			Websocket: HTTPWsConfig{
				Enabled:          false,
				ListenPort:       8669,
				ListenAddress:    "0.0.0.0",
				CrossOriginValue: "*",
			},
			Socket: DomainSocket{
				Enabled: true,
				Path:    "gocc.ipc",
			},
		},
		P2P: P2P{
			MaxPeers:           20,
			ListenPort:         10209,
			ListenAddress:      "localhost",
			ConnectionTimeout:  40,
			MinPeersThreashold: 2,
			Bootstraper: Bootstraper{
				BootstrapPeriodic: 120,
			},
		},
	}
}

// LoadTomlConfig loads the toml file into the config structure
func LoadTomlConfig(ctx *cli.Context, cfg *GlobalConfig) {
	conf := ctx.String("config")
	if conf == "" {
		log.Fatal("Configuration file not given")
	}
	fh, err := os.Open(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	if err := toml.NewDecoder(fh).Decode(cfg); err != nil {
		log.Fatal(err)
	}
}

// ApplyFlags applies the flags and sets the globalconfig
func ApplyFlags(ctx *cli.Context, cfg *GlobalConfig) {
	// Global
	if ctx.GlobalIsSet(LogLevelFlag.Name) {
		cfg.Global.LogLevel = ctx.GlobalString(LogLevelFlag.Name)
	}
	if ctx.GlobalIsSet(DataDirFlag.Name) {
		cfg.Global.DataDir = ctx.GlobalString(DataDirFlag.Name)
	}
	if ctx.GlobalIsSet(KeystoreDirFlag.Name) {
		cfg.Global.KeystoreDir = ctx.GlobalString(KeystoreDirFlag.Name)
	}
	if ctx.GlobalIsSet(UploadsDirFlag.Name) {
		cfg.Global.UploadsDir = ctx.GlobalString(UploadsDirFlag.Name)
	}
	if ctx.GlobalIsSet(DatabaseNameFlag.Name) {
		cfg.Global.DatabaseName = ctx.GlobalString(DatabaseNameFlag.Name)
	}
	if ctx.GlobalIsSet(AvailabilityFlag.Name) {
		cfg.Global.Availability = strings.Split(ctx.GlobalString(AvailabilityFlag.Name), ",")
	}

	// Host
	if ctx.GlobalIsSet(MaxContainersFlag.Name) {
		cfg.Host.MaxContainers = ctx.GlobalInt(MaxContainersFlag.Name)
	}
	if ctx.GlobalIsSet(CPUPerContainerFlag.Name) {
		cfg.Host.CPUPerContainer = ctx.GlobalInt(CPUPerContainerFlag.Name)
	}

	if ctx.GlobalIsSet(GPUPerContainerFlag.Name) {
		cfg.Host.GPUPerContainer = ctx.GlobalInt(GPUPerContainerFlag.Name)
	}
	if ctx.GlobalIsSet(MemoryPerContainerFlag.Name) {
		cfg.Host.MemoryPerContainer = ctx.GlobalInt(MemoryPerContainerFlag.Name)
	}

	if ctx.GlobalIsSet(StoragePerContainerFlag.Name) {
		cfg.Host.MemoryPerContainer = ctx.GlobalInt(StoragePerContainerFlag.Name)
	}
	if ctx.GlobalIsSet(DockerSwarmAdvertiseAddrFlag.Name) {
		cfg.Host.DockerSwarm.AdvertiseAddress = ctx.GlobalString(DockerSwarmAdvertiseAddrFlag.Name)
	}
	if ctx.GlobalIsSet(DockerSwarmAddrFlag.Name) {
		cfg.Host.DockerSwarm.ListenAddress = ctx.GlobalString(DockerSwarmAddrFlag.Name)
	}
	if ctx.GlobalIsSet(DockerSwarmPortFlag.Name) {
		cfg.Host.DockerSwarm.ListenPort = ctx.GlobalInt(DockerSwarmPortFlag.Name)
	}

	// RPC
	if ctx.GlobalIsSet(RPCFlag.Name) {
		cfg.RPC.Enabled = ctx.GlobalBool(RPCFlag.Name)
	}

	if ctx.GlobalIsSet(RPCServicesFlag.Name) {
		cfg.RPC.EnabledServices = strings.Split(ctx.GlobalString(RPCServicesFlag.Name), ",")
	}

	if ctx.GlobalIsSet(RPCWhitelistFlag.Name) {
		cfg.RPC.Whitelist = strings.Split(ctx.GlobalString(RPCWhitelistFlag.Name), ",")
	}

	// RPC:SOCKET
	if ctx.GlobalIsSet(RPCSocketEnabledFlag.Name) {
		cfg.RPC.Socket.Enabled = ctx.GlobalBool(RPCSocketEnabledFlag.Name)
	}

	if ctx.GlobalIsSet(RPCSocketPathFlag.Name) {
		cfg.RPC.Socket.Path = ctx.GlobalString(RPCSocketPathFlag.Name)
	}

	// RPC:HTTP
	if ctx.GlobalIsSet(RPCHTTPEnabledFlag.Name) {
		cfg.RPC.HTTP.Enabled = ctx.GlobalBool(RPCHTTPEnabledFlag.Name)
	}

	if ctx.GlobalIsSet(RPCHTTPPortFlag.Name) {
		cfg.RPC.HTTP.ListenPort = ctx.GlobalInt(RPCHTTPPortFlag.Name)
	}

	if ctx.GlobalIsSet(RPCHTTPAddrFlag.Name) {
		cfg.RPC.HTTP.ListenAddress = ctx.GlobalString(RPCHTTPAddrFlag.Name)
	}

	if ctx.GlobalIsSet(RPCHTTPCrossOriginFlag.Name) {
		cfg.RPC.HTTP.CrossOriginValue = ctx.GlobalString(RPCHTTPCrossOriginFlag.Name)
	}

	// RPC:WS
	if ctx.GlobalIsSet(RPCWSEnabledFlag.Name) {
		cfg.RPC.Websocket.Enabled = ctx.GlobalBool(RPCWSEnabledFlag.Name)
	}

	if ctx.GlobalIsSet(RPCWSPortFlag.Name) {
		cfg.RPC.Websocket.ListenPort = ctx.GlobalInt(RPCWSPortFlag.Name)
	}

	if ctx.GlobalIsSet(RPCWSAddrFlag.Name) {
		cfg.RPC.Websocket.ListenAddress = ctx.GlobalString(RPCWSAddrFlag.Name)
	}

	if ctx.GlobalIsSet(RPCWSCrossOriginFlag.Name) {
		cfg.RPC.Websocket.CrossOriginValue = ctx.GlobalString(RPCWSCrossOriginFlag.Name)
	}

	// P2P
	if ctx.GlobalIsSet(MaxPeersFlag.Name) {
		cfg.P2P.MaxPeers = ctx.GlobalInt(MaxPeersFlag.Name)
	}

	if ctx.GlobalIsSet(P2PListenPortFlag.Name) {
		cfg.P2P.ListenPort = ctx.GlobalInt(P2PListenPortFlag.Name)
	}

	if ctx.GlobalIsSet(P2PListenAddrFlag.Name) {
		cfg.P2P.ListenAddress = ctx.GlobalString(P2PListenAddrFlag.Name)
	}

	if ctx.GlobalIsSet(P2PTimeoutFlag.Name) {
		cfg.P2P.ConnectionTimeout = ctx.GlobalInt(P2PTimeoutFlag.Name)
	}

	if ctx.GlobalIsSet(P2PMinPeerThreasholdFlag.Name) {
		cfg.P2P.MinPeersThreashold = ctx.GlobalInt(P2PMinPeerThreasholdFlag.Name)
	}

	if ctx.GlobalIsSet(P2PBootstraperFlag.Name) {
		cfg.P2P.Bootstraper.Nodes = strings.Split(ctx.GlobalString(P2PBootstraperFlag.Name), ",")
	}

	if ctx.GlobalIsSet(P2PPeriodicFlag.Name) {
		cfg.P2P.Bootstraper.BootstrapPeriodic = ctx.GlobalInt(P2PPeriodicFlag.Name)
	}

}

// GetConfig returns the config after applying the defaults, toml file and flags
func GetConfig(ctx *cli.Context) *GlobalConfig {
	// create default config
	cfg := DefaultConfig()
	// if config file is given, load it
	confFile := ctx.String("config")
	if confFile != "" {
		LoadTomlConfig(ctx, cfg)
	}
	// apply flags to config
	ApplyFlags(ctx, cfg)
	return cfg
}

// DefaultDataDir is the default data directory to use for the databases and other
// persistence requirements.
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "GoCC_data")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "GoCC_data")
		} else {
			return filepath.Join(home, ".gocc_data")
		}
	}
	// if a common location wasn't found, return empty and handle later
	return ""
}

// Return the user's home dir
func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
