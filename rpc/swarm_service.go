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
	"github.com/docker/docker/api/types/swarm"
)

// SwarmService used to register Docker image functionality
// over jsonrpc
type SwarmService struct{}

// Init initiates a swarm manager
func (s *SwarmService) Init(ctx context.Context, advertise string, listen string) (string, error) {
	res, err := manager.GetInstance().InitSwarm(advertise, listen)
	if err != nil {
		return res, err
	}
	return res, nil
}

// Leave the swarm mode
func (s *SwarmService) Leave(ctx context.Context) (bool, error) {
	res, err := manager.GetInstance().LeaveSwarm()
	if err != nil {
		return res, err
	}
	return res, nil
}

// Join the swarm mode
func (s *SwarmService) Join(ctx context.Context, advertise string, datapath string, remoteAddrs []string, token string, listen string) (bool, error) {
	res, err := manager.GetInstance().SwarmJoin(advertise, datapath, remoteAddrs, token, listen)
	if err != nil {
		return res, err
	}
	return res, nil
}

// Inspect the swarm
func (s *SwarmService) Inspect(ctx context.Context) (swarm.Swarm, error) {
	swrm, err := manager.GetInstance().SwarmInspect()
	if err != nil {
		return swrm, err
	}
	return swrm, nil
}

// ServiceCreate creates a swarm service
func (s *SwarmService) ServiceCreate(ctx context.Context, service swarm.ServiceSpec, options types.ServiceCreateOptions) (types.ServiceCreateResponse, error) {
	resp, err := manager.GetInstance().ServiceCreate(service, options)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
