package rpc

import (
	"context"

	"github.com/crowdcompute/crowdengine/manager"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// ContainerService used to register Docker container functionality
type ContainerService struct{}

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
