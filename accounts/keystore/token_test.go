package keystore

import (
	"testing"

	"github.com/crowdcompute/crowdengine/crypto"
	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {
	symmKey, err := crypto.RandomEntropy(32)
	if err != nil {
		t.Errorf("There was an error getting random entrop. Error: %s", err)
	}
	token, err := NewToken(symmKey, NewTokenClaims("", ""))
	if err != nil {
		t.Errorf("There was an error creating new token. Error: %s", err)
	}
	verification, err := VerifyToken(token, symmKey)
	if err != nil {
		t.Errorf("There was an error verifying token for account. Error: %s", err)
	}
	assert.True(t, verification)
}
