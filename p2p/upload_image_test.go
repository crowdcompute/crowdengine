package p2p

import (
	"testing"
)

type tester struct {
	input  interface{}
	expect bool
}

var (
	tests = []tester{
		{"image ID does not exist here", false},
		{"sha256:imageID exists here", true},
		{"", false},
	}
)

func TestGetImageID(t *testing.T) {
	for _, test := range tests {
		if _, exists := getImageID(test.input.(string)); exists != test.expect {
			t.Errorf("test = %v, expext = %v: It is %v that image ID in this test text!", test.input, test.expect, exists)
		}
	}
}
