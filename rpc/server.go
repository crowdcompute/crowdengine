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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/crowdcompute/crowdengine/accounts/keystore"
	"github.com/crowdcompute/crowdengine/common"
	"github.com/crowdcompute/crowdengine/log"

	"github.com/ethereum/go-ethereum/rpc"
)

func ServeHTTP(apis []API, ks *keystore.KeyStore) http.HandlerFunc {
	server := rpc.NewServer()
	for _, api := range apis {
		err := server.RegisterName(api.Namespace, api.Service)
		common.FatalIfErr(err, "Ethereum RPC could not register name.")
	}
	return authRequired(apis, ks, server)
}

// authRequired is a middleware for the HTTP server.
// Authenticates a token and passes the request to the next handler
func authRequired(apis []API, ks *keystore.KeyStore, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if empty body
		if r.ContentLength == 0 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		protected, err := isMethodProtected(apis, buf.Bytes())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// Restore the r.Body to its original state
		r.Body = ioutil.NopCloser(buf)

		// ns is protected, place the logic which verifies the header
		if protected {
			key, err := getKeyForAccount(ks, r.Header)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), common.ContextKeyPair, key)
			log.Printf("Token valid and account {%s} unlocked. ", key.Address)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		next.ServeHTTP(w, r)
	}
}

func isMethodProtected(apis []API, rawJSONBody []byte) (bool, error) {
	namespace, method, err := FindNamespaceMethod(rawJSONBody)
	if err != nil {
		return false, err
	}
	// find which namespace
	namespaceMethodProtected := false
	for _, v := range apis {
		if v.Namespace == namespace {
			// if * then all methods are protected
			if v.AuthRequired == "*" {
				namespaceMethodProtected = true
				break
			}

			// break them and inspect them
			fncs := strings.Split(v.AuthRequired, ",")
			for _, w := range fncs {
				if common.LcFirst(strings.TrimSpace(w)) == method {
					namespaceMethodProtected = true
					break
				}
			}
			break
		}
	}
	return namespaceMethodProtected, nil
}

// Extracts the token from authorization header,
// and checks if token valid and related acount unlocked.
// And returns the key
func getKeyForAccount(ks *keystore.KeyStore, header http.Header) (*keystore.Key, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		err := fmt.Errorf("No Authorization given on header")
		log.Println(err.Error())
		return nil, err
	}
	token := strings.Split(authHeader, " ")[1]
	key, err := ks.GetKeyIfUnlockedAndValid(token)
	if err != nil {
		log.Println("Error while trying to get key for a token. Error: ", err)
		return nil, err
	}
	return key, nil
}
