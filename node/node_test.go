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

package node

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/common"
	ccrpc "github.com/crowdcompute/crowdengine/rpc"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

var (
	listenPort    = 10209
	listenAddress = "127.0.0.1"
	testNode, _   = NewNode(&config.GlobalConfig{
		RPC: config.RPC{Enabled: true, HTTP: config.HTTPWsConfig{Enabled: true, ListenPort: listenPort, ListenAddress: listenAddress}},
	})
)

func createPOSTreq(jsonStr []byte) *http.Request {
	url := fmt.Sprintf("http://%s:%d/", listenAddress, listenPort)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func serveRequest(n *Node, req *http.Request) *httptest.ResponseRecorder {
	server := rpc.NewServer()
	for _, api := range testApis() {
		err := server.RegisterName(api.Namespace, api.Service)
		common.FatalIfErr(err, "Ethereum RPC could not register name.")
	}
	handler := authRequired(testApis(), n.ks, server)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func testApis() []ccrpc.API {
	return []ccrpc.API{
		{
			Namespace:    "testService",
			Version:      "1.0",
			Service:      &TestService{},
			Public:       true,
			AuthRequired: "MethodThatRequiresAuth",
		},
	}
}

type TestService struct{}

func (api *TestService) MethodThatRequiresAuth(ctx context.Context) error {
	return nil
}

func (api *TestService) MethodThatNotRequiresAuth(ctx context.Context) error {
	return nil
}

func TestMethodRequiresAuthButNotAuthGiven(t *testing.T) {
	var jsonStr = []byte(`{"jsonrpc":"2.0","id":"1","method":"testService_methodThatRequiresAuth","params":[]}`)
	req := createPOSTreq(jsonStr)
	resp := serveRequest(testNode, req)

	assert.True(t, resp.Code == http.StatusUnauthorized)
}

func TestMethodNotRequiresAuth(t *testing.T) {
	var jsonStr = []byte(`{"jsonrpc":"2.0","id":"1","method":"testService_methodThatNotRequiresAuth","params":[]}`)
	req := createPOSTreq(jsonStr)
	resp := serveRequest(testNode, req)
	assert.True(t, resp.Code == http.StatusOK)
}

func TestMethodNotGiven(t *testing.T) {
	var jsonStr = []byte(`{"jsonrpc":"2.0","id":"1","method":"testService","params":[]}`)
	req := createPOSTreq(jsonStr)
	resp := serveRequest(testNode, req)
	assert.True(t, resp.Code == http.StatusBadRequest)
}
