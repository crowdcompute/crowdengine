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
	"bufio"
	"fmt"
	"time"

	"github.com/crowdcompute/crowdengine/log"

	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/manager"
	api "github.com/crowdcompute/crowdengine/p2p/protomsgs"

	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
	protobufCodec "github.com/multiformats/go-multicodec/protobuf"
	uuid "github.com/satori/go.uuid"
)

// pattern: /protocol-name/request-or-response-message/version
const runRequest = "/task/availreq/0.0.1"
const runResponse = "/task/availresp/0.0.1"

// This struct implements the Notifier interface
// TaskProtocol type
type TaskProtocol struct {
	p2pHost           host.Host // local host
	ContainerID       chan string
	runningContainers []string
	taskObservers     map[Observer]struct{}
}

func NewTaskProtocol(p2pHost host.Host) *TaskProtocol {
	p := &TaskProtocol{p2pHost: p2pHost,
		ContainerID:       make(chan string, 1),
		runningContainers: make([]string, 0),
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
	// get request data
	data := &api.RunRequest{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	common.CheckErr(err, "[onRunRequest] Couldn't decode data.")

	log.Printf("%s: Received avail request from %s.", s.Conn().LocalPeer(), s.Conn().RemotePeer())

	valid := authenticateProtoMsg(data, data.RunImageMsgData.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	containerID := createRunContainer(data.ImageID)

	// generate response message
	log.Printf("%s: Sending run image response to %s. Message id: %s...", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.RunImageMsgData.MessageData.Id)

	resp := &api.RunResponse{RunImageMsgData: NewRunImageMsgData(data.RunImageMsgData.MessageData.Id, false, p.p2pHost),
		ContainerID: containerID}

	key := p.p2pHost.Peerstore().PrivKey(p.p2pHost.ID())
	resp.RunImageMsgData.MessageData.Sign = signProtoMsg(resp, key)

	// send the response
	if sendMsg(p.p2pHost, s.Conn().RemotePeer(), resp, protocol.ID(runResponse)) {
		log.Printf("%s: Run image response to %s sent.", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	}

	p.runningContainers = append(p.runningContainers, containerID)
	log.Println("Start tracking job's status...")
	// Start tracking jobs' status
	// Wait here for the job to finish.
	// TODO: We shouldn't wait here, we should have a switch
	go p.waitForJobToFinish(containerID)
}

// Running until job's done
func (p *TaskProtocol) waitForJobToFinish(containerID string) {
	log.Println("start task status tracking")
	// TODO: Time has to be a const somewhere
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Checking if job's done...")
			if !containerRunning(containerID) {
				deleteValFromSlice(p.runningContainers, containerID)
				log.Println("Job's done checking pending requests...")
				p.Notify()
				return
			}
		}
	}
}

func deleteValFromSlice(slice []string, val string) []string {
	var newSlice []string
	for _, v := range slice {
		if v == val {
			continue
		} else {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}

func createRunContainer(imageID string) string {
	container, err := manager.GetInstance().CreateContainer(imageID)
	_, err = manager.GetInstance().RunContainer(container.ID)
	common.CheckErr(err, fmt.Sprintf("Error running the container %s: %s", imageID, err))
	return container.ID
}

// remote ping response handler
func (p *TaskProtocol) onRunResponse(s inet.Stream) {
	data := &api.RunResponse{}
	decoder := protobufCodec.Multicodec(nil).Decoder(bufio.NewReader(s))
	err := decoder.Decode(data)
	common.CheckErr(err, "[onRunResponse] Couldn't decode data.")

	valid := authenticateProtoMsg(data, data.RunImageMsgData.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	log.Printf("%s: Received running image response from %s. Message id:%s.", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.RunImageMsgData.MessageData.Id)
	p.ContainerID <- data.ContainerID
}
