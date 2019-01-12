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
