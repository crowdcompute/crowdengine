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

// Inspect the swarm
func (s *SwarmService) ServiceCreate(ctx context.Context, service swarm.ServiceSpec, options types.ServiceCreateOptions) (types.ServiceCreateResponse, error) {
	resp, err := manager.GetInstance().ServiceCreate(service, options)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
