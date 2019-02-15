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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/crowdcompute/crowdengine/accounts/keystore"
	"github.com/crowdcompute/crowdengine/cmd/gocc/config"
	"github.com/crowdcompute/crowdengine/common"
	ccrpc "github.com/crowdcompute/crowdengine/rpc"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

var (
	tempDir       = ""
	passphrase    = "testPass"
	listenPort    = 10209
	listenAddress = "127.0.0.1"
	testNode      *Node
)

func init() {
	var err error
	tempDir, err = ioutil.TempDir("", "ccompute-keystore-test")
	common.FatalIfErr(err, "Couldn't create temp file for node_test.go")
	testNode, _ = NewNode(&config.GlobalConfig{
		RPC:    config.RPC{Enabled: true, HTTP: config.HTTPWsConfig{Enabled: true, ListenPort: listenPort, ListenAddress: listenAddress}},
		Global: config.Global{KeystoreDir: tempDir},
	})
}

func createPOSTreq(jsonStr []byte) *http.Request {
	url := fmt.Sprintf("http://%s:%d/", listenAddress, listenPort)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func serveRequest(ks *keystore.KeyStore, req *http.Request) *httptest.ResponseRecorder {
	server := rpc.NewServer()
	for _, api := range testApis() {
		err := server.RegisterName(api.Namespace, api.Service)
		common.FatalIfErr(err, "Ethereum RPC could not register name.")
	}
	handler := authRequired(testApis(), ks, server)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func testApis() []ccrpc.API {
	return []ccrpc.API{
		{
			Namespace:    "testService",
			Version:      "1.0",
			Service:      &TestService{ks: testNode.ks},
			Public:       true,
			AuthRequired: "MethodThatRequiresAuth",
		},
	}
}

type TestService struct {
	ks *keystore.KeyStore
}

func (api *TestService) CreateTestAccount(ctx context.Context, passphrase string) (string, error) {
	acc, err := api.ks.NewAccount(passphrase)
	common.FatalIfErr(err, "There was an error creating the account")
	return acc.Address, err
}

func (api *TestService) UnlockTestAccount(ctx context.Context, accAddr, passphrase string) (string, error) {
	// First issue a token
	rawToken, err := api.ks.IssueTokenForAccount(accAddr, keystore.NewTokenClaims("", ""))
	if err != nil {
		return "", err
	}
	// Then unlock the account if there is no issue with the Token creation above
	if err := api.ks.TimedUnlock(accAddr, passphrase, common.TokenTimeout); err != nil {
		return "", err
	}
	return rawToken, err
}

func (api *TestService) CreateTTAccount(ctx context.Context, passphrase string) error {
	return nil
}

func (api *TestService) MethodThatRequiresAuth(ctx context.Context) error {
	return nil
}

func (api *TestService) MethodThatNotRequiresAuth(ctx context.Context) error {
	return nil
}

func TestMethodRequiresAuthButNotAuthGiven(t *testing.T) {
	var jsonStr = []byte(`{"jsonrpc":"2.0","id":"1","method":"testService_methodThatRequiresAuth","params":[]}`)
	req := createPOSTreq(jsonStr)
	resp := serveRequest(testNode.ks, req)

	assert.True(t, resp.Code == http.StatusUnauthorized)
}

type EtherRPCResponse struct {
	Result string `json:"result"`
}

// This is an integration test.
// Creates an account, unlocks it, and passes it as Authirization to the header on an RPC method that requires Authirization
func TestMethodRequiresAuthAndAuthGiven(t *testing.T) {
	// Create a temporary account using json RPC http request
	var jsonStr = []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":"1","method":"testService_createTestAccount","params":["%s"]}`, passphrase))
	req := createPOSTreq(jsonStr)
	resp := serveRequest(testNode.ks, req)
	fmt.Println("resp.Body.String()")
	fmt.Println(resp.Body.String())
	defer os.RemoveAll(tempDir)

	respjson := EtherRPCResponse{}
	err := json.Unmarshal(resp.Body.Bytes(), &respjson)
	if err != nil {
		t.Errorf("There was an error unmarshaling Ethereum RPC response. Their response might have changed!")
	}
	accAddr := respjson.Result
	fmt.Println(accAddr)

	// Unlock account
	var jsonUnlockAcc = []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":"1","method":"testService_unlockTestAccount","params":["%s","%s"]}`, accAddr, passphrase))
	reqUnlockAcc := createPOSTreq(jsonUnlockAcc)
	respUnlockAcc := serveRequest(testNode.ks, reqUnlockAcc)
	fmt.Println("respUnlockAcc.Body.String()")
	fmt.Println(respUnlockAcc.Body.String())
	respUnlockJSON := EtherRPCResponse{}
	err = json.Unmarshal(respUnlockAcc.Body.Bytes(), &respUnlockJSON)
	if err != nil {
		t.Errorf("There was an error unmarshaling Ethereum RPC response. Their response might have changed!")
	}
	token := respUnlockJSON.Result
	fmt.Println(token)

	// Check if token valid
	var jsonReqAuth = []byte(`{"jsonrpc":"2.0","id":"1","method":"testService_methodThatRequiresAuth","params":[]}`)
	reqRequiresAuth := createPOSTreq(jsonReqAuth)
	reqRequiresAuth.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	respRequiresAuth := serveRequest(testNode.ks, reqRequiresAuth)
	assert.True(t, respRequiresAuth.Code == http.StatusOK)
}

func TestMethodNotRequiresAuth(t *testing.T) {
	var jsonStr = []byte(`{"jsonrpc":"2.0","id":"1","method":"testService_methodThatNotRequiresAuth","params":[]}`)
	req := createPOSTreq(jsonStr)
	resp := serveRequest(testNode.ks, req)
	assert.True(t, resp.Code == http.StatusOK)
}

func TestMethodNotGiven(t *testing.T) {
	var jsonStr = []byte(`{"jsonrpc":"2.0","id":"1","method":"testService","params":[]}`)
	req := createPOSTreq(jsonStr)
	resp := serveRequest(testNode.ks, req)
	assert.True(t, resp.Code == http.StatusBadRequest)
}
