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
	"log"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/manager"
	"github.com/crowdcompute/crowdengine/p2p"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
)

type ServiceAPI struct {
	host *p2p.Host
}

type Task struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
}

func NewServiceAPI(h *p2p.Host) *ServiceAPI {
	return &ServiceAPI{host: h}
}

func (s *ServiceAPI) Run(ctx context.Context, task string) error {
	t := Task{}
	json.Unmarshal([]byte(task), &t)
	fmt.Println("I got this task:", t)

	// Initialize a docker Swarm
	_, err := manager.GetInstance().InitSwarm(s.host.IP, "0.0.0.0:2377")
	common.CheckErr(err, "[onUploadResponse] Couldn't initialize swarm.")

	if swarmInspect, err := manager.GetInstance().SwarmInspect(); err == nil {
		// TODO: Check if user wants a Manager or Worker. Some nodes might be managers
		s.host.WorkerToken = swarmInspect.JoinTokens.Worker
		s.host.ManagerToken = swarmInspect.JoinTokens.Manager
	} else {
		log.Printf("Error doing Swarm Inspect: %v", err)
		return err
	}

	// Send Join request to node's bootnodes
	s.host.SendJoinToNeighbours(t.Replicas)
	service, err := s.createService(&t)

	if err != nil {
		fmt.Printf("Cannot create service. Error: %s", err)
		return err
	}

	log.Printf("Service created successfully! %s\n", service)

	return nil
}

func (s *ServiceAPI) createService(task *Task) (string, error) {
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
	serviceCreate, err := manager.GetInstance().ServiceCreate(serviceSpec, types.ServiceCreateOptions{})

	if err != nil {
		fmt.Printf("Cannot create service. Error: %s", err)
		return "", err
	}

	return serviceCreate.ID, nil
}
