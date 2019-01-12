package hexutil

import "testing"

type marshalTest struct {
	input interface{}
	want  string
}

var (
	encodeBytesToHexTests = []marshalTest{
		{[]byte{}, "0x"},
		{[]byte{0}, "0x00"},
		{[]byte{0, 'a'}, "0x0061"},
	}
)

func TestEncode(t *testing.T) {
	for _, test := range encodeBytesToHexTests {
		enc := Encode(test.input.([]byte))
		if enc != test.want {
			t.Errorf("input %x: wrong encoding %s", test.input, enc)
		}
	}
}
