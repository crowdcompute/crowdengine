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

package crypto

import (
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	keypair, err := GenerateKeyPair()
	if err != nil {
		t.Errorf("Error while generating assymetric key: %s", err)
	}
	privBytes, _ := keypair.Private.Bytes()
	if len(privBytes) != 36 {
		t.Errorf("Public key size is incorrect, got: %d, want: %d.", len(privBytes), 36)
	}
	if len(keypair.Address) != 40 {
		t.Errorf("Address size is incorrect, got: %d, want: %d.", len(keypair.Address), 40)
	}
}
