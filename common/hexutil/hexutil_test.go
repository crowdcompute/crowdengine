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

package hexutil

import "testing"

type tester struct {
	input  interface{}
	expect string
}

var (
	encodeBytesToHexTests = []tester{
		{[]byte{}, "0x"},
		{[]byte{0}, "0x00"},
		{[]byte{0, 'a'}, "0x0061"},
	}
)

func TestEncode(t *testing.T) {
	for _, test := range encodeBytesToHexTests {
		enc := EncodeWithPrefix(test.input.([]byte))
		if enc != test.expect {
			t.Errorf("Input %x: Wrong hex encoding %s", test.input, enc)
		}
	}
}
