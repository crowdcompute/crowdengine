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

package keystore

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalUnmarshal(t *testing.T) {
	key := NewKey()
	dummyPass := "test"
	data, err := key.MarshalJSON(dummyPass)
	if err != nil {
		t.Errorf("There was an error marshaling the key: %s", err)
	}
	unmarshaledKey, err := UnmarshalKey(data, dummyPass)
	if err != nil {
		t.Errorf("There was an error unmarshaling the key: %s", err)
	}
	priv, err := key.Private.Bytes()
	if err != nil {
		t.Errorf("There was an error getting the bytes from the priv key: %s", err)
	}
	unmarshaledPriv, err := unmarshaledKey.Private.Bytes()
	if err != nil {
		t.Errorf("There was an error getting the bytes from the priv key: %s", err)
	}
	assert.True(t, hex.EncodeToString(priv) == hex.EncodeToString(unmarshaledPriv))
}
