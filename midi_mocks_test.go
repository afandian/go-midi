// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

// Tests for test mocks.

package midi

import (
	"testing"
)

/* Tests for tests */

// The MockReadSeeker should behaves as a ReadSeeker should.
func TestMockReadSeekerReads(t *testing.T) {
	var reader = NewMockReadSeeker(&[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07})

	// Buffer to read into.
	var data []byte = []byte{0x00, 0x00, 0x00}
	var count = 0

	// Start with empty data buffer
	assertBytesEqual(data, []byte{0x00, 0x00, 0x00}, t)

	// Now read into the 3-length buffer
	count, err := reader.Read(data)
	if count != 3 {
		t.Fatal("Count not 3 was ", count)
	}

	if err != nil {
		t.Fatal("Error not nil, was ", err)
	}

	assertBytesEqual(data, []byte{0x01, 0x02, 0x03}, t)

	// Read into it again to get the next 3
	count, err = reader.Read(data)
	if count != 3 {
		t.Fatal("Count not 3 was ", count)
	}

	if err != nil {
		t.Fatal("Error not nil, was ", err)
	}
	assertBytesEqual(data, []byte{0x04, 0x05, 0x06}, t)

	// Read again to get the last one.
	count, err = reader.Read(data)
	if count != 1 {
		t.Fatal("Count not 1 was ", count)
	}

	if err != nil {
		t.Fatal("Error not nil, was ", err)
	}

	// Data will still have the old data remaining
	assertBytesEqual(data, []byte{0x07, 0x05, 0x06}, t)

	// One more time, should be empty
	count, err = reader.Read(data)
	if count != 0 {
		t.Fatal("Count not 0 was ", count)
	}

	if err != nil {
		t.Fatal("Error not nil, was ", err)
	}
}

// The MockReadSeeker should behaves as a ReadSeeker should.
func TestMockReadSeekerSeeks(t *testing.T) {
	var reader = NewMockReadSeeker(&[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07})
	var dataSize int64 = 7

	var count = 0

	// Single byte buffer.
	var byteBuffer []byte = []byte{0x00}

	/*
	 * 0 - Relative to start of file 
	 */

	// Seek from the start of the file to the last byte.
	sook, err := reader.Seek(dataSize-1, 0)

	if err != nil {
		t.Fatal(err)
	}

	if sook != dataSize-1 {
		t.Fatal("Expected to return ", dataSize-1, " got ", sook)
	}

	count, err = reader.Read(byteBuffer)

	if err != nil {
		t.Fatal(err)
	}

	if count != 1 {
		t.Fatal("Expected to read 1 byte, got ", count)
	}

	if byteBuffer[0] != 0x07 {
		t.Fatal("Expected 0x07 got ", byteBuffer)
	}

	// Seek from the start of the file to the 3rd byte.
	sook, err = reader.Seek(2, 0)

	if err != nil {
		t.Fatal(err)
	}

	if sook != 2 {
		t.Fatal("Expected to return ", 2, " got ", sook)
	}

	count, err = reader.Read(byteBuffer)

	if err != nil {
		t.Fatal(err)
	}

	if byteBuffer[0] != 0x03 {
		t.Fatal("Expected 0x03 got ", byteBuffer)
	}

	// Seek from the start of the file to the first byte.
	sook, err = reader.Seek(0, 0)

	if err != nil {
		t.Fatal(err)
	}

	if sook != 0 {
		t.Fatal("Expected to return ", 0, " got ", sook)
	}

	count, err = reader.Read(byteBuffer)

	if err != nil {
		t.Fatal(err)
	}

	if byteBuffer[0] != 0x01 {
		t.Fatal("Expected 0x01 got ", byteBuffer)
	}

	/*
	 * 1 - Relative to current position
	 */

	// Seek from the current position to the same place.

	// Get in the middle
	sook, err = reader.Seek(4, 0)

	if err != nil {
		t.Fatal(err)
	}

	if sook != 4 {
		t.Fatal("Expected to return ", 4, " got ", sook)
	}

	// Seek same place relative to current.
	sook, err = reader.Seek(0, 1)

	if err != nil {
		t.Fatal(err)
	}

	if sook != 4 {
		t.Fatal("Expected to return ", 4, " got ", sook)
	}

	count, err = reader.Read(byteBuffer)

	if err != nil {
		t.Fatal(err)
	}

	if byteBuffer[0] != 0x05 {
		t.Fatal("Expected 0x05 got ", byteBuffer)
	}

	// Seek forward a byte
	sook, err = reader.Seek(1, 1)

	if err != nil {
		t.Fatal(err)
	}

	if sook != 6 {
		t.Fatal("Expected to return ", 6, " got ", sook)
	}

	count, err = reader.Read(byteBuffer)

	if err != nil {
		t.Fatal(err)
	}

	if byteBuffer[0] != 0x07 {
		t.Fatal("Expected 0x07 got ", byteBuffer)
	}

	/*
	 * 2 - Relative to end of file
	 */

	// Seek from the current position to the same place.

	// Get to the end.
	// Seek same place relative to end
	sook, err = reader.Seek(0, 2)

	if err != nil {
		t.Fatal(err)
	}

	if sook != 7 {
		t.Fatal("Expected to return ", 7, " got ", sook)
	}

	count, err = reader.Read(byteBuffer)

	if err != nil {
		t.Fatal(err)
	}

	if byteBuffer[0] != 0x07 {
		t.Fatal("Expected 0x07 got ", byteBuffer)
	}

	// Seek back a byte
	sook, err = reader.Seek(1, 2)

	if err != nil {
		t.Fatal(err)
	}

	if sook != 6 {
		t.Fatal("Expected to return ", 6, " got ", sook)
	}

	count, err = reader.Read(byteBuffer)

	if err != nil {
		t.Fatal(err)
	}

	if byteBuffer[0] != 0x07 {
		t.Fatal("Expected 0x07 got ", byteBuffer)
	}
}
