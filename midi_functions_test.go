// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

// Tests for test_functions. Work directly with value parsing out.

package midi

import (
	"io"
	"testing"
)

// Test that parseVarLength does what it should.
// Example data taken from http://www.music.mcgill.ca/~ich/classes/mumt306/midiformat.pdf
func TestVarLengthParser(t *testing.T) {
	expected := []uint32{
		0x00000000,
		0x00000040,
		0x0000007F,
		0x00000080,
		0x00002000,
		0x00003FFF,
		0x00004000,
		0x00100000,
		0x001FFFFF,
		0x00200000,
		0x08000000,
		0x0FFFFFFF}

	input := [][]byte{
		[]byte{0x00},
		[]byte{0x40},
		[]byte{0x7F},
		[]byte{0x81, 0x00},
		[]byte{0xC0, 0x00},
		[]byte{0xFF, 0x7F},
		[]byte{0x81, 0x80, 0x00},
		[]byte{0xC0, 0x80, 0x00},
		[]byte{0xFF, 0xFF, 0x7F},
		[]byte{0x81, 0x80, 0x80, 0x00},
		[]byte{0xC0, 0x80, 0x80, 0x00},
		[]byte{0xFF, 0xFF, 0xFF, 0x7F}}

	var numItems = len(input)

	for i := 0; i < numItems; i++ {
		var reader = NewMockReadSeeker(&input[i])
		var result, err = parseVarLength(reader)

		if result != expected[i] {
			t.Fatal("Expected ", expected[i], " got ", result)
		}

		if err != nil {
			t.Fatal("Expected no error got ", err)
		}
	}

	// Now we want to read past the end of a var length file.
	// We should get an UnexpectedEndOfFile error.

	// First read OK.
	var reader = NewMockReadSeeker(&input[0])
	var _, err = parseVarLength(reader)
	if err != nil {
		t.Fatal("Expected no error got ", err)
	}

	// Second read not OK.
	_, err = parseVarLength(reader)
	if err != UnexpectedEndOfFile {
		t.Fatal("Expected End of file ")
	}
}

// Test for parseUint32
func TestParse32Bit(t *testing.T) {
	expected := []uint32{
		0,
		1,
		4,
		42,
		429,
		4294,
		42949,
		429496,
		4294967,
		42949672,
		429496729,
		4294967295,
	}

	input := [][]byte{
		[]byte{0x00, 0x00, 0x00, 0x00},
		[]byte{0x00, 0x00, 0x00, 0x01},
		[]byte{0x00, 0x00, 0x00, 0x04},
		[]byte{0x00, 0x00, 0x00, 0x2A},
		[]byte{0x00, 0x00, 0x01, 0xAD},
		[]byte{0x00, 0x00, 0x10, 0xC6},
		[]byte{0x00, 0x00, 0xA7, 0xC5},
		[]byte{0x00, 0x06, 0x8D, 0xB8},
		[]byte{0x00, 0x41, 0x89, 0x37},
		[]byte{0x02, 0x8F, 0x5C, 0x28},
		[]byte{0x19, 0x99, 0x99, 0x99},
		[]byte{0xFF, 0xFF, 0xFF, 0xFF},
	}

	var numItems = len(input)

	for i := 0; i < numItems; i++ {
		var reader = NewMockReadSeeker(&input[i])
		var result, err = parseUint32(reader)

		if result != expected[i] {
			t.Fatal("Expected ", expected[i], " got ", result)
		}

		if err != nil {
			t.Fatal("Expected no error got ", err)
		}
	}

	// Now we want to read past the end of a file.
	// We should get an UnexpectedEndOfFile error.

	// First read OK.
	var reader = NewMockReadSeeker(&input[0])
	var _, err = parseUint32(reader)
	if err != nil {
		t.Fatal("Expected no error got ", err)
	}

	// Second read not OK.
	_, err = parseUint32(reader)
	if err != UnexpectedEndOfFile {
		t.Fatal("Expected End of file ")
	}
}

// Test for parseUint16
func TestParse16Bit(t *testing.T) {
	expected := []uint16{
		0,
		1,
		4,
		42,
		429,
		4294,
		42949,
	}

	input := [][]byte{
		[]byte{0x00, 0x00},
		[]byte{0x00, 0x01},
		[]byte{0x00, 0x04},
		[]byte{0x00, 0x2A},
		[]byte{0x01, 0xAD},
		[]byte{0x10, 0xC6},
		[]byte{0xA7, 0xC5},
	}

	var numItems = len(input)

	for i := 0; i < numItems; i++ {
		var reader = NewMockReadSeeker(&input[i])
		var result, err = parseUint16(reader)

		if result != expected[i] {
			t.Fatal("Expected ", expected[i], " got ", result)
		}

		if err != nil {
			t.Fatal("Expected no error got ", err)
		}
	}

	// Now we want to read past the end of a file.
	// We should get an UnexpectedEndOfFile error.

	// First read OK.
	var reader = NewMockReadSeeker(&input[0])
	var _, err = parseUint16(reader)
	if err != nil {
		t.Fatal("Expected no error got ", err)
	}

	// Second read not OK.
	_, err = parseUint16(reader)
	if err != UnexpectedEndOfFile {
		t.Fatal("Expected End of file ")
	}
}

// Test for parseChunkHeader.
func TestParseChunkHeader(t *testing.T) {
	// Headers for two popular chunk types.

	// Length 6, as all MThds shoudl be 6 long.
	var MThd = []byte{0x4D, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06}

	// Arbitrary length 4294967 (parseUint32 is tested separately).
	var MTrk = []byte{0x4D, 0x54, 0x72, 0x6b, 0x00, 0x41, 0x89, 0x37}

	// Too short in the type word.
	var tooShort1 = []byte{0x4D, 0x54, 0x72}

	// To short in the length word.
	var tooShort2 = []byte{0x4D, 0x54, 0x72, 0x6b, 0x00, 0x41, 0x89}

	var header ChunkHeader
	var err error
	var reader io.ReadSeeker

	// Try for typical MIDI file header
	reader = NewMockReadSeeker(&MThd)
	header, err = parseChunkHeader(reader)

	if header.chunkType != "MThd" {
		t.Fatal("Got ", header, " expected MThd")
	}

	if header.length != 6 {
		t.Fatal("Got ", header, " expected 6")
	}

	if err != nil {
		t.Fatal("Got error ", err)
	}

	// Try for typical track header
	reader = NewMockReadSeeker(&MTrk)
	header, err = parseChunkHeader(reader)

	if header.chunkType != "MTrk" {
		t.Fatal("Got ", header, " expected MTrk")
	}

	if header.length != 4294967 {
		t.Fatal("Got ", header, " expected 4294967")
	}

	if err != nil {
		t.Fatal("Got error ", err)
	}

	// Now two incomplete headers.

	// Too short to parse the type
	reader = NewMockReadSeeker(&tooShort1)
	header, err = parseChunkHeader(reader)

	if err == nil {
		t.Fatal("Expected error for tooshort1")
	}

	// Too short to parse the length
	reader = NewMockReadSeeker(&tooShort2)
	header, err = parseChunkHeader(reader)

	if err == nil {
		t.Fatal("Expected error for tooshort 2")
	}
}

// Test for parseHeaderData.
func TestParseHeaderData(t *testing.T) {
	var err error
	var data, expected HeaderData

	// Format: 1
	// Tracks: 2
	// Division: metrical 5
	var headerMetrical = NewMockReadSeeker(&[]byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x05})
	expected = HeaderData{
		format:              1,
		numTracks:           2,
		timeFormat:          MetricalTimeFormat,
		timeFormatData:      0x00,
		ticksPerQuarterNote: 5}

	data, err = parseHeaderData(headerMetrical)

	if err != nil {
		t.Fatal(err)
	}

	if data != expected {
		t.Fatal(data, " != ", expected)
	}

	// Format: 2
	// Tracks: 1
	// Division: timecode (actual data ignored for now)
	var headerTimecode = NewMockReadSeeker(&[]byte{0x00, 0x02, 0x00, 0x01, 0xFF, 0x05})
	expected = HeaderData{
		format:              2,
		numTracks:           1,
		timeFormat:          TimeCodeTimeFormat,
		timeFormatData:      0x7F05, // Removed the top timecode type bit flag.
		ticksPerQuarterNote: 0}

	data, err = parseHeaderData(headerTimecode)

	if err != nil {
		t.Fatal(err)
	}

	if data != expected {
		t.Fatal(data, " != ", expected)
	}

	// Format: 3, which doesn't exist.
	var badFormat = NewMockReadSeeker(&[]byte{0x00, 0x03, 0x00, 0x01, 0xFF, 0x05})
	data, err = parseHeaderData(badFormat)

	if err != UnsupportedSmfFormat {
		t.Fatal("Expected exception but got none")
	}

	// Too short in each field
	var tooShort1 = NewMockReadSeeker(&[]byte{0x00, 0x02, 0x00, 0x01, 0xFF})
	data, err = parseHeaderData(tooShort1)

	if err != UnexpectedEndOfFile {
		t.Fatal("Expected exception but got ", err)
	}

	var tooShort2 = NewMockReadSeeker(&[]byte{0x00, 0x02, 0x00})
	data, err = parseHeaderData(tooShort2)

	if err != UnexpectedEndOfFile {
		t.Fatal("Expected exception but got none")
	}

	var tooShort3 = NewMockReadSeeker(&[]byte{0x00})
	data, err = parseHeaderData(tooShort3)

	if err != UnexpectedEndOfFile {
		t.Fatal("Expected exception but got none")
	}
}
