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

package p2p

import (
	"time"

	"github.com/crowdcompute/crowdengine/common"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/manager"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"

	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	uuid "github.com/satori/go.uuid"
)

// pattern: /protocol-name/request-or-response-message/version
const runRequest = "/task/availreq/0.0.1"
const runResponse = "/task/availresp/0.0.1"

// TaskProtocol implements the Notifier interface
type TaskProtocol struct {
	p2pHost           host.Host // local host
	ContainerID       chan string
	runningContainers map[string]struct{} // The running containers of the node
	taskObservers     map[Observer]struct{}
}

// NewTaskProtocol sets the protocol's stream handlers and returns a new TaskProtocol
func NewTaskProtocol(p2pHost host.Host) *TaskProtocol {
	p := &TaskProtocol{p2pHost: p2pHost,
		ContainerID:       make(chan string, 1),
		runningContainers: map[string]struct{}{},
		taskObservers:     map[Observer]struct{}{},
	}
	p2pHost.SetStreamHandler(runRequest, p.onRunRequest)
	p2pHost.SetStreamHandler(runResponse, p.onRunResponse)
	return p
}

// Register an observer to get notified when a job is done
func (p *TaskProtocol) Register(o Observer) {
	p.taskObservers[o] = struct{}{}
}

// Deregister an observer on runtime
func (p *TaskProtocol) Deregister(o Observer) {
	delete(p.taskObservers, o)
}

// Notify all registered observers
func (p *TaskProtocol) Notify() {
	for observer := range p.taskObservers {
		observer.onNotify()
	}
}

// RunImage runs an image with imageID to the hostID
func (p *TaskProtocol) RunImage(hostID peer.ID, imageID string) bool {
	log.Printf("%s: Asking running image. Sending request to: %s....", p.p2pHost.ID(), hostID)
	// create message data
	req := &api.RunRequest{RunImageMsgData: NewRunImageMsgData(uuid.Must(uuid.NewV4(), nil).String(), true, p.p2pHost),
		ImageID: imageID}

	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	req.RunImageMsgData.MessageData.Sign = signProtoMsg(req, key)

	if !sendMsg(p.p2pHost, hostID, req, protocol.ID(runRequest)) {
		return false
	}

	log.Printf("%s: Ask running image to: %s was sent. Message Id: %s", p.p2pHost.ID(), peer.ID(hostID), req.RunImageMsgData.MessageData.Id)
	return true
}

// remote peer requests handler
func (p *TaskProtocol) onRunRequest(s inet.Stream) {
	log.Printf("%s: Received run container request from %s.", s.Conn().LocalPeer(), s.Conn().RemotePeer())
	// get request data
	data := &api.RunRequest{}
	decodeProtoMessage(data, s)

	if valid := authenticateProtoMsg(data, data.RunImageMsgData.MessageData); !valid {
		log.Println("Failed to authenticate message")
		return
	}
	containerID, err := manager.GetInstance().CreateRunContainer(data.ImageID)
	if err != nil {
		log.Errorf("Error crating a container. Error: %s", err)
		return
	}
	p.createSendResponse(s.Conn().RemotePeer(), containerID)
	p.runningContainers[containerID] = struct{}{}
	log.Println("Start tracking job's status...")

	go p.waitForJobToFinish(containerID)
}

// Create and send a response to the toPeer note
func (p *TaskProtocol) createSendResponse(toPeer peer.ID, response string) bool {
	log.Printf("%s: Sending run image response to %s.", p.p2pHost.ID(), toPeer)

	resp := &api.RunResponse{RunImageMsgData: NewRunImageMsgData(uuid.Must(uuid.NewV4(), nil).String(), false, p.p2pHost),
		ContainerID: response}

	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.RunImageMsgData.MessageData.Sign = signProtoMsg(resp, key)

	// send the response
	sentOK := sendMsg(p.p2pHost, toPeer, resp, protocol.ID(runResponse))
	if sentOK {
		log.Printf("%s: Run image response to %s was sent.", p.p2pHost.ID(), toPeer)
	}
	return sentOK
}

// Start tracking jobs' status
func (p *TaskProtocol) waitForJobToFinish(containerID string) {
	ticker := time.NewTicker(common.ContainerCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Checking if job's done...")
			if !containerRunning(containerID) {
				delete(p.runningContainers, containerID)
				log.Println("Job's done checking pending requests...")
				p.Notify()
				return
			}
		}
	}
}

// remote ping response handler
func (p *TaskProtocol) onRunResponse(s inet.Stream) {
	data := &api.RunResponse{}
	decodeProtoMessage(data, s)

	valid := authenticateProtoMsg(data, data.RunImageMsgData.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	log.Printf("%s: Received running image response from %s. Message id:%s.", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.RunImageMsgData.MessageData.Id)
	p.ContainerID <- data.ContainerID
}
