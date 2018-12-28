package manager

import "github.com/docker/docker/client"

// Image represents an image
type Image struct {
	ID   string
	Name string
}

// Container represents a container
type Container struct {
	ID   string
	Name string
}

// Manager represents the container manager
type DockerManager struct {
	client     *client.Client
	Images     []Image
	Containers []Container
}
