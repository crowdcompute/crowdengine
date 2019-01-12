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

package common

import "testing"

type tester struct {
	input1 interface{}
	input2 interface{}
	expect bool
}

type testStruct1 struct {
	testProperty int
}
type testStruct2 struct {
	testProperty1 int
	testProperty2 string
}

var (
	sliceExistsTests = []tester{
		{[]testStruct1{{1}, {2}}, testStruct1{1}, true},
		{[]testStruct1{{1}, {2}}, testStruct1{2}, true},
		{[]testStruct1{{1}, {2}}, testStruct1{3}, false},

		{[]testStruct2{{1, ""}, {2, ""}}, testStruct1{1}, false},
		{[]testStruct2{{1, "test"}, {2, "test"}}, testStruct2{2, "test"}, true},
		{[]testStruct2{{1, "test"}, {2, "test"}}, testStruct2{2, "wrong test"}, false},
	}
)

func TestSliceExists(t *testing.T) {
	for _, test := range sliceExistsTests {
		exists := SliceExists(test.input1, test.input2)
		if exists != test.expect {
			t.Errorf("slice %v, item: %v: If this item should exist in slice, the answer should be %v", test.input1, test.input2, exists)
		}
	}
}

func TestPanic(t *testing.T) {
	assertPanic(t, FunctionThatPanics)
}

func FunctionThatPanics() {
	// First parameter is not a slice so the function will panic
	SliceExists(testStruct1{}, testStruct1{})
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
