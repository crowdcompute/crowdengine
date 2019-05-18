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

package manager

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"sync"
	"fmt"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

var (
	instance *DockerManager
	once     sync.Once
)

// GetInstance gets single manager - Thread-safe
func GetInstance() *DockerManager {
	once.Do(func() {
		var containers []Container
		var images []Image
		mngr := &DockerManager{Images: images, Containers: containers}
		if mngr.boot() {
			instance = mngr
		}
	})
	return instance
}

// boot the manager
func (m *DockerManager) boot() bool {
	cli, err := client.NewEnvClient()

	if err != nil {
		return false
	}
	m.client = cli
	return true
}

// BuildImageFromDockerfile builds an image from a Dockerfile
func (m *DockerManager) BuildImageFromDockerfile() bool {
	return false
}

// LoadImage loads a complete image
func (m *DockerManager) LoadImage(filePath string) (string, error) {
	destination, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	response, err := m.client.ImageLoad(context.Background(), destination, true)
	defer response.Body.Close()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	r, _ := regexp.Compile("Loaded image(.*)")
	matches := r.FindAllStringSubmatch(string(body), -1)
	// If there is no match then the file is not right
	if len(matches) == 0 {
		return "", errors.New("Can't load image, the file is wrong")
	}
	result := matches[0][1]
	return result, nil
}

// ListImages list all the available images
func (m *DockerManager) ListImages(options types.ImageListOptions) ([]types.ImageSummary, error) {
	images, err := m.client.ImageList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	return images, nil
}

// Pulling an image from a registry
func (m *DockerManager) ImagePull(refStr string, options types.ImagePullOptions) (string, error) {
	r, err := m.client.ImagePull(context.Background(), refStr, options)
	if err != nil {
		log.Fatalf("Error while calling docker api's ImagePull: %s", err)
		return "", err
	}
	var b []byte
	r.Read(b)

	return string(b), nil
}

// RemoveImage an image from docker
func (m *DockerManager) RemoveImage(imageID string, options types.ImageRemoveOptions) ([]types.ImageDelete, error) {
	r, err := m.client.ImageRemove(context.Background(), imageID, options)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// ListContainers list all the container
// Shows all running containers by docker daemon
// It might include other container not registered by Crowdcompute
func (m *DockerManager) ListContainers() ([]types.Container, error) {
	containers, err := m.client.ContainerList(context.Background(), types.ContainerListOptions{All: true, Latest: true})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// ListRegisteredContainers list all the container regitered by this node
func (m *DockerManager) ListRegisteredContainers() ([]types.Container, error) {
	var registered []types.Container
	containers, err := m.ListContainers()
	if err != nil {
		return registered, err
	}

	for _, v := range containers {
		for _, t := range m.Containers {
			if t.ID == v.ID {
				registered = append(registered, v)
			}
		}
	}
	return registered, nil
}

// Logs shows logs of a specific container
func (m *DockerManager) Logs(containerid string, sincetime string) ([]byte, error) {
	log, err := m.client.ContainerLogs(context.Background(), containerid, types.ContainerLogsOptions{Since: sincetime, ShowStdout: true, ShowStderr: true, Timestamps: true})
	defer log.Close()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(log)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func newVolumeMount(src, dst string) mount.Mount {
	return mount.Mount{
		Type:         mount.TypeVolume,
		Source:       src,
		Target:       dst,
		ReadOnly:     false,
		BindOptions:  nil,
		TmpfsOptions: nil,
		VolumeOptions: &mount.VolumeOptions{
			NoCopy: false,
			Labels: map[string]string{},
		},
	}
}

// CreateContainer the manager
// TODO: persist containerid into levelDB
func (m *DockerManager) CreateContainer(imageID string) (container.ContainerCreateCreatedBody, error) {
	ctx := context.Background()
	hostconfig := new(container.HostConfig)
	hostconfig.Mounts = make([]mount.Mount, 0)
	hostconfig.Mounts = append(hostconfig.Mounts, newVolumeMount(imageID, common.DockerMountDest)) // imageID will be the name of the volume
	// TODO: Give permissions to edit the /home folder
	resp, err := m.client.ContainerCreate(ctx, &container.Config{
		Image: imageID,
	}, hostconfig, nil, "")

	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}
	m.Containers = append(m.Containers, Container{ID: resp.ID})
	return resp, nil
}

// RunContainer the manager
func (m *DockerManager) RunContainer(containerid string) (bool, error) {
	ctx := context.Background()
	if err := m.client.ContainerStart(ctx, containerid, types.ContainerStartOptions{}); err != nil {
		return false, err
	}
	return true, nil
}

// CreateRunContainer creates and runs a container from an image ID
func (m *DockerManager) CreateRunContainer(imageID string) (string, error) {
	container, err := m.CreateContainer(imageID)
	if err != nil {
		return "", fmt.Errorf("Error creating container form this image ID: %s. Image ID could be wrong. Error: %s", imageID, err)
	}
	_, err = m.RunContainer(container.ID)
	return container.ID, err
}


// InspectContainer inspects a running container
func (m *DockerManager) InspectContainer(containerid string) (types.ContainerJSON, error) {
	inspection, err := m.client.ContainerInspect(context.Background(), containerid)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return inspection, nil
}

// InspectContainerRaw inspects a running container
func (m *DockerManager) InspectContainerRaw(containerid string, getSize bool) (types.ContainerJSON, []byte, error) {
	inspection, raw, err := m.client.ContainerInspectWithRaw(context.Background(), containerid, getSize)
	if err != nil {
		return types.ContainerJSON{}, nil, err
	}
	return inspection, raw, nil
}

// RemoveContainer removes a container
func (m *DockerManager) RemoveContainer(containerid string, options types.ContainerRemoveOptions) error {
	err := m.client.ContainerRemove(context.Background(), containerid, options)
	if err != nil {
		return err
	}
	return nil
}
