package keystore

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalUnmarshal(t *testing.T) {
	key := NewKey()
	dummyPass := "test"
	data, err := MarshalKey(dummyPass, key)
	if err != nil {
		t.Errorf("There was an error marshaling the key: %s", err)
	}
	unmarshaledKey, err := UnmarshalKey(data, dummyPass)
	if err != nil {
		t.Errorf("There was an error unmarshaling the key: %s", err)
	}
	priv, _ := key.Private.Bytes()
	unmarshaledPriv, _ := unmarshaledKey.Private.Bytes()
	assert.True(t, hex.EncodeToString(priv) == hex.EncodeToString(unmarshaledPriv))
}
