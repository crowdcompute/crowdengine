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
	// RPCAddrFlag rpc host and port to connect to
	RPCAddrFlag = cli.StringFlag{
		Name:  "rpcaddr",
		Usage: "HTTP-RPC host and port to connect to",
	}

	// FileserverFlag host and port to connect to
	FileserverFlag = cli.StringFlag{
		Name:  "fileserver",
		Usage: "Fileserver host and port to connect to",
	}

	// AccAddrFlag is the account address to lock or unlock
	AccAddrFlag = cli.StringFlag{
		Name:  "account",
		Usage: "Account address to lock or unlock",
	}

	// AccPassphraseFlag is the passphrase needed to unlock the account
	AccPassphraseFlag = cli.StringFlag{
		Name:  "passphrase",
		Usage: "Passphrase to unlock an account",
	}

	// ImgPathFlag is the filepath to the image to be uploaded
	ImgPathFlag = cli.StringFlag{
		Name:  "imgpath",
		Usage: "imgpath to the image to be uploaded",
	}

	// Libp2pIDFlag is the libp2pid of the node
	Libp2pIDFlag = cli.StringFlag{
		Name:  "libp2pid",
		Usage: "libp2pid of the node",
	}

	// ImgIDFlag is the docker image id to run on the node
	ImgIDFlag = cli.StringFlag{
		Name:  "imgid",
		Usage: "docker image id to run on the node",
	}

	// ServiceNameFlag is the docker swarm name to give to the service 
	ServiceNameFlag = cli.StringFlag{
		Name:  "servicename",
		Usage: "docker swarm name to give to the service ",
	}

	// ServiceImgFlag is the docker swarm image to run as a service
	ServiceImgFlag = cli.StringFlag{
		Name:  "serviceimg",
		Usage: "docker swarm image to run as a service",
	}
)
