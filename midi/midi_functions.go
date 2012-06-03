// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

// Functions for reading MIDI data.

package midi

import (
	"io"
)

// parseUint32 parse a 4-byte 32 bit integer from a Reader.
func parseUint32(reader io.Reader) (uint32, error) {

	var buffer []byte = make([]byte, 4)
	num, err := reader.Read(buffer)

	// If we couldn't read 4 bytes, that's a problem.
	if num != 4 {
		return 0, UnexpectedEndOfFile
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return 0, err
	}

	var value uint32 = 0x00
	value |= uint32(buffer[3]) << 0
	value |= uint32(buffer[2]) << 8
	value |= uint32(buffer[1]) << 16
	value |= uint32(buffer[0]) << 24

	return value, nil
}

// parseVarLength parse a variable length value from a Reader.
func parseVarLength(reader io.Reader) (uint32, error) {

	// Single byte buffer to read byte by byte.
	var buffer []byte = make([]uint8, 1)

	// The number of bytes returned.
	// Should always be 1 unless we reach the EOF
	var num int = 1

	// Result value
	var result uint32 = 0x00

	// RTFM.
	var first = true
	for (first || (buffer[0] & 0x80 == 0x80)) && (num > 0) {
		result = result << 7

		num, _ = reader.Read(buffer)
		result |= (uint32(buffer[0]) & 0x7f)
		first = false
	}

	if num == 0 && ! first {
		return result, UnexpectedEndOfFile
	}

	return result, nil
}

// parseChunkHeader parses a chunk header from a Reader.
func parseChunkHeader(reader io.Reader) (*ChunkHeader, error) {
	var chunk ChunkHeader

	var chunkTypeBuffer []byte = make([]byte, 4)
	num, err := reader.Read(chunkTypeBuffer)

	// If we couldn't read 4 bytes, that's a problem.
	if num != 4 {
		return &chunk, UnexpectedEndOfFile
	}

	if err != nil {
		return &chunk, err
	}

	chunk.length, err = parseUint32(reader)
	chunk.chunkType = string(chunkTypeBuffer)

	// parseUint32 might return an error.
	if err != nil {
		return &chunk, err
	}

	return &chunk, nil
}