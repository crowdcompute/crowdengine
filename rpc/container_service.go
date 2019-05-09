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

package rpc

import (
	"context"

	"github.com/crowdcompute/crowdengine/manager"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// ContainerService used to register Docker container functionality
type ContainerService struct{}

// NewContainerService returns a new ContainerService
func NewContainerService() *ContainerService {
	return &ContainerService{}
}

// RawResult represents any result
type RawResult struct {
	data interface{}
}

// List list all the available images
func (s *ContainerService) List(ctx context.Context) ([]types.Container, error) {
	containers, err := manager.GetInstance().ListContainers()
	if err != nil {
		return containers, err
	}
	return containers, nil
}

// ListRegistered lists the registered containers by this
func (s *ContainerService) ListRegistered(ctx context.Context) ([]types.Container, error) {
	containers, err := manager.GetInstance().ListRegisteredContainers()
	if err != nil {
		return containers, err
	}
	return containers, nil
}

// Logs prints logs of a container
func (s *ContainerService) Logs(ctx context.Context, containerid string, sincetime string) ([]manager.DockerLog, error) {
	logs, err := manager.GetInstance().Logs(containerid, sincetime)
	if err != nil {
		return nil, err
	}
	decoded, err := manager.DockerLogDecoder(logs)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

// Create creates a container from an image
func (s *ContainerService) Create(ctx context.Context, imageid string) (container.ContainerCreateCreatedBody, error) {
	container, err := manager.GetInstance().CreateContainer(imageid)
	if err != nil {
		return container, err
	}
	return container, nil
}

// Run runs a container
func (s *ContainerService) Run(ctx context.Context, containerid string) (bool, error) {
	running, err := manager.GetInstance().RunContainer(containerid)
	if err != nil {
		return running, err
	}
	return running, nil
}

// Inspect inspects a container
func (s *ContainerService) Inspect(ctx context.Context, containerid string) (types.ContainerJSON, error) {
	containerStatus, err := manager.GetInstance().InspectContainer(containerid)
	if err != nil {
		return containerStatus, err
	}
	return containerStatus, nil
}
