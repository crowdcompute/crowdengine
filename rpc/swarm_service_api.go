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
	"encoding/json"
	"fmt"

	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/manager"
	"github.com/crowdcompute/crowdengine/p2p"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
)

// SwarmServiceAPI ...
type SwarmServiceAPI struct {
	host *p2p.Host
	cfg  *config.DockerSwarm
}

type swarmTask struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

// NewSwarmServiceAPI creates a new RPC service with methods specific for creating a swarm service.
func NewSwarmServiceAPI(h *p2p.Host, cfg *config.DockerSwarm) *SwarmServiceAPI {
	return &SwarmServiceAPI{host: h, cfg: cfg}
}

// Run initializes a swarm, makes nodes to join it, and creates a swarm service
func (api *SwarmServiceAPI) Run(ctx context.Context, task string, nodes []string) error {
	t := swarmTask{}
	json.Unmarshal([]byte(task), &t)
	log.Println("I got this task:", t)

	if err := api.initSwarm(); err != nil {
		return err
	}
	api.host.SendJoinToPeersAndWait(nodes)

	serviceID, err := createService(&t)
	if err != nil {
		log.Printf("Cannot create service. Error: %s", err)
		return err
	}
	log.Printf("Service created successfully! %s\n", serviceID)
	return nil
}

// initSwarm initializes a docker Swarm and stores the swarm's worker & manager
// tokens in memory
func (api *SwarmServiceAPI) initSwarm() error {
	swarmListen := fmt.Sprintf("%s:%d", api.cfg.ListenAddress, api.cfg.ListenPort)
	_, err := manager.GetInstance().InitSwarm(api.host.IP, swarmListen)
	if err != nil {
		return err
	}
	if swarmInspect, errInspect := manager.GetInstance().SwarmInspect(); errInspect == nil {
		// TODO: Check if user has to be a Manager or Worker. Some nodes might be managers
		api.host.WorkerToken = swarmInspect.JoinTokens.Worker
		api.host.ManagerToken = swarmInspect.JoinTokens.Manager
	} else {
		log.Printf("Error running Swarm Inspect: %v", errInspect)
		return errInspect
	}
	return err
}

// createService creates and starts a swarm service
func createService(task *swarmTask) (string, error) {
	serviceSpec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: task.Name,
			Labels: map[string]string{
				"key1": "",
			},
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: task.Image,
			},
		},
		// Mode: swarm.ServiceMode{},
	}
	serviceCreateRes, err := manager.GetInstance().ServiceCreate(serviceSpec, types.ServiceCreateOptions{})
	if err != nil {
		log.Printf("Cannot create service. Error: %s", err)
		return "", err
	}
	return serviceCreateRes.ID, nil
}

// Stop makes all nodes involved to leave the swarm
func (api *SwarmServiceAPI) Stop(ctx context.Context, nodes []string) error {
	if _, err := manager.GetInstance().LeaveSwarm(); err != nil {
		return err
	}
	api.host.SendLeaveToPeersAndWait(nodes)
	return nil
}
