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

package common

import "time"

// FileChunk is the size of a chunk when uploading a file
const FileChunk = 1 * (1 << 20)

// ImagesDest is the destination folder for storing images
const ImagesDest = "./uploads/"

// TenDays represents 10 days in time
const TenDays time.Duration = 24 * time.Hour * 10

// TTLmsg represents the time to live of a p2p message
const TTLmsg time.Duration = time.Second * 15

// FillChar is a character that fills a chunck of bytes
// ':' it's an illegal character for file names under windows and linux
const FillChar string = ":"

// FileSizeLength represents the total length of the buffer whenever sending the file size to a peer
const FileSizeLength int = 10

// FileNameLength represents the total length of the buffer whenever sending the file name to a peer
const FileNameLength int = 64

// SignatureLength represents the total length of the buffer whenever sending the signature to a peer
// TODO: Not sure what number to give here. Need to see the range
const SignatureLength int = 150

// HashLength represents the total length of the buffer whenever sending the hash to a peer
// TODO: Not sure what number to give here. Need to see the range
const HashLength int = 100

// RemoveImagesInterval represents the time interval to check for removing images
const RemoveImagesInterval time.Duration = time.Second * 10

// ContainerCheckInterval represents the time interval to check whether a container has finished running
const ContainerCheckInterval time.Duration = time.Second * 3

// DiscoveryTimeout represents the time to wait for
const DiscoveryTimeout time.Duration = time.Second * 10

// TokenTimeout represents the time to wait for the JWT account token to expire
const TokenTimeout time.Duration = time.Second * 60

type ContextKey string

// ContextKeyPair represents the context key name for a private key
const ContextKeyPair ContextKey = "keypair"

// ContextKeyUploadPath represents the context key name for an upload path
const ContextKeyUploadPath ContextKey = "uploadpath"

// DockerMountDest is the destination of the mount of all docker containers
const DockerMountDest string = "/home/data"
