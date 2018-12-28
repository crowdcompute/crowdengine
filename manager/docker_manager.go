package manager

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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
	// Resolving the below error:
	// Error response from daemon: client version 1.39 is too new. Maximum supported API version is 1.37
	//cli, err := client.NewClientWithOpts(client.WithVersion("1.37"))

	cli, err := client.NewEnvClient()

	if err != nil {
		return false
	}
	m.client = cli
	return true
}

//************************************************************************//
//**************************** IMAGES ************************************//
//************************************************************************//

// BuildImageFromDockerfile builds an image from a Dockerfile
func (m *DockerManager) BuildImageFromDockerfile() bool {
	return false
}

// LoadImage loads a complete image
func (m *DockerManager) LoadImage(filename string) (string, error) {
	destination, err := os.Open(common.ImagesDest + filename)
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

//************************************************************************//
//**************************** CONTAINERS ********************************//
//************************************************************************//

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

// CreateContainer the manager
// TODO: persist containerid into levelDB
func (m *DockerManager) CreateContainer(imageid string) (container.ContainerCreateCreatedBody, error) {
	ctx := context.Background()
	resp, err := m.client.ContainerCreate(ctx, &container.Config{
		Image: imageid,
	}, nil, nil, "")

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

// InspectContainer inspects a running container
func (m *DockerManager) InspectContainer(containerid string) (types.ContainerJSON, error) {
	inspection, err := m.client.ContainerInspect(context.Background(), containerid)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return inspection, nil
}

// InspectContainer inspects a running container
func (m *DockerManager) InspectContainerRaw(containerid string, getSize bool) (types.ContainerJSON, []byte, error) {
	inspection, raw, err := m.client.ContainerInspectWithRaw(context.Background(), containerid, getSize)
	if err != nil {
		return types.ContainerJSON{}, nil, err
	}
	return inspection, raw, nil
}

// InspectContainer inspects a running container
func (m *DockerManager) RemoveContainer(containerid string, options types.ContainerRemoveOptions) error {
	err := m.client.ContainerRemove(context.Background(), containerid, options)
	if err != nil {
		return err
	}
	return nil
}
