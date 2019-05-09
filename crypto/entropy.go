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
	"crypto/rand"
	"errors"
	"io"
)

// RandomEntropy returns a slice of n bytes from rand.Reader
func RandomEntropy(length int) ([]byte, error) {
	buf := make([]byte, length)
	n, err := io.ReadFull(rand.Reader, buf)
	if err != nil || n != length {
		return nil, errors.New("Unable to read random bytes of length")
	}
	return buf, nil
}
