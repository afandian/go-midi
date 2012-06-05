// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

// Test helper functions.

package midi

import (
	"testing"
)

// Helper assertions
func assertHasFlag(value int, flag int, test *testing.T) {
	if value&flag == 0 {
		test.Fatal("Expected to find ", flag, " in ", value)
	}
}

// assertBytesEqual asserts that the given byte arrays or slices are equal.
func assertBytesEqual(a []byte, b []byte, t *testing.T) {
	if len(a) != len(b) {
		t.Fatal("Two arrays not equal", a, b)
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Fatal("Two arrays not equal. At ", i, " was ", a[i], " vs ", b[i])
		}
	}
}

// Assert uint16s equal
func assertUint16Equal(a int, b int, test *testing.T) {
	if a != b {
		test.Fatal(a, " != ", b)
	}
}
