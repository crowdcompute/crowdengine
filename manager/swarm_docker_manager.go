package manager

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
)

// InitSwarm initiates a swarm manager
func (m *DockerManager) InitSwarm(advertise string, listen string) (string, error) {
	swarmconfig := swarm.InitRequest{AdvertiseAddr: advertise, ListenAddr: listen}
	str, err := m.client.SwarmInit(context.Background(), swarmconfig)
	if err != nil {
		return "", err
	}
	return str, nil
}

// LeaveSwarm leaves a swarm
func (m *DockerManager) LeaveSwarm() (bool, error) {
	err := m.client.SwarmLeave(context.Background(), true)
	if err != nil {
		return false, err
	}
	return true, nil
}

// SwarmJoin joins a swarm
func (m *DockerManager) SwarmJoin(advertise string, datapath string, remoteAddrs []string, token string, listen string) (bool, error) {
	joinconfig := swarm.JoinRequest{AdvertiseAddr: advertise, RemoteAddrs: remoteAddrs, JoinToken: token, ListenAddr: listen}
	err := m.client.SwarmJoin(context.Background(), joinconfig)
	if err != nil {
		return false, err
	}
	return true, nil
}

// SwarmInspect inspects a swarm
func (m *DockerManager) SwarmInspect() (swarm.Swarm, error) {
	swrm, err := m.client.SwarmInspect(context.Background())
	if err != nil {
		return swrm, err
	}
	return swrm, nil
}

// SwarmInfo returns info about the swarm
func (m *DockerManager) SwarmInfo() (swarm.Info, error) {
	info, err := m.client.Info(context.Background())
	if err != nil {
		return info.Swarm, err
	}
	return info.Swarm, nil
}

// ServiceCreate creates a docker swarm service
func (m *DockerManager) ServiceCreate(service swarm.ServiceSpec, options types.ServiceCreateOptions) (types.ServiceCreateResponse, error) {
	resp, err := m.client.ServiceCreate(context.Background(), service, options)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
