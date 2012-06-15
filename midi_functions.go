// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

/*
 * Functions for reading actual MIDI data in the various formats that crop up.
 */

package midi

import (
	"fmt"
	"io"
)

// parseUint32 parse a 4-byte 32 bit integer from a ReadSeeker.
// It returns the 32-bit value and an error.
func parseUint32(reader io.ReadSeeker) (uint32, error) {
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

// parseUint16 parses a 2-byte 16 bit integer from a ReadSeeker.
// It returns the 16-bit value and an error.
func parseUint16(reader io.ReadSeeker) (uint16, error) {

	var buffer []byte = make([]byte, 2)
	num, err := reader.Read(buffer)

	// If we couldn't read 2 bytes, that's a problem.
	if num != 2 {
		return 0, UnexpectedEndOfFile
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return 0, err
	}

	var value uint16 = 0x00
	value |= uint16(buffer[1]) << 0
	value |= uint16(buffer[0]) << 8

	return value, nil
}

// parseUint7 parses a 7-bit bit integer from a ReadSeeker, ignoring the high bit.
// It returns the 8-bit value and an error.
func parseUint7(reader io.ReadSeeker) (uint8, error) {

	var buffer []byte = make([]byte, 1)
	num, err := reader.Read(buffer)

	// If we couldn't read 1 bytes, that's a problem.
	if num != 1 {
		return 0, UnexpectedEndOfFile
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return 0, err
	}

	return (buffer[0] & 0x7f), nil
}

// parseTwoUint7 parses two 7-bit bit integer stored in consecutive bytes from a ReadSeeker, ignoring the high bit in each.
// It returns the 8-bit value and an error.
func parseTwoUint7(reader io.ReadSeeker) (uint8, uint8, error) {

	var buffer []byte = make([]byte, 2)
	num, err := reader.Read(buffer)

	// If we couldn't read 2 bytes, that's a problem.
	if num != 2 {
		return 0, 0, UnexpectedEndOfFile
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return 0, 0, err
	}

	return (buffer[0] & 0x7f), (buffer[1] & 0x7f), nil
}

// parseUint8 parses an 8-bit bit integer stored in a bytes from a ReadSeeker.
// It returns a single uint8.
func parseUint8(reader io.ReadSeeker) (uint8, error) {

	var buffer []byte = make([]byte, 1)
	num, err := reader.Read(buffer)

	// If we couldn't read 1 bytes, that's a problem.
	if num != 1 {
		return 0, UnexpectedEndOfFile
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return 0, err
	}

	return uint8(buffer[0]), nil
}


// parsePitchWheelValue parses a 14-bit signed value, which becomes a signed int16.
// The least significant 7 bits are stored in the first byte, the 7 most significant bites are stored in the second.
// Return the signed value relative to the centre, and the absolute value.
// This is tested in midi_lexer_test.go TestPitchWheel
func parsePitchWheelValue(reader io.ReadSeeker) (relative int16, absolute uint16, err error) {

	var buffer []byte = make([]byte, 2)
	num, err := reader.Read(buffer)

	// If we couldn't read 2 bytes, that's a problem.
	if num != 2 {
		return 0, 0, UnexpectedEndOfFile
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return 0, 0, err
	}

	var val uint16 = 0

	val = uint16((buffer[1])&0x7f) << 7
	val |= uint16(buffer[0]) & 0x7f
	fmt.Println(val)

	// log.Println()
	// Turn into a signed value relative to the centre.
	relative = int16(val) - 0x2000

	return relative, val, nil
}

// parseVarLength parses a variable length value from a ReadSeeker.
// It returns the [up to] 32-bit value and an error.
func parseVarLength(reader io.ReadSeeker) (uint32, error) {

	// Single byte buffer to read byte by byte.
	var buffer []byte = make([]uint8, 1)

	// The number of bytes returned.
	// Should always be 1 unless we reach the EOF
	var num int = 1

	// Result value
	var result uint32 = 0x00

	// RTFM.
	var first = true
	for (first || (buffer[0]&0x80 == 0x80)) && (num > 0) {
		result = result << 7

		num, _ = reader.Read(buffer)
		result |= (uint32(buffer[0]) & 0x7f)
		first = false
	}

	if num == 0 && !first {
		return result, UnexpectedEndOfFile
	}

	return result, nil
}

// parseChunkHeader parses a chunk header from a ReadSeeker.
// It returns the ChunkHeader struct as a value and an error.
func parseChunkHeader(reader io.ReadSeeker) (ChunkHeader, error) {
	var chunk ChunkHeader

	var chunkTypeBuffer []byte = make([]byte, 4)
	num, err := reader.Read(chunkTypeBuffer)

	// If we couldn't read 4 bytes, that's a problem.
	if num != 4 {
		return chunk, UnexpectedEndOfFile
	}

	if err != nil {
		return chunk, err
	}

	chunk.length, err = parseUint32(reader)
	chunk.chunkType = string(chunkTypeBuffer)

	// parseUint32 might return an error.
	if err != nil {
		return chunk, err
	}

	return chunk, nil
}

// parseHeaderData parses SMF-header chunk header data.
// It returns the ChunkHeader struct as a value and an error.
func parseHeaderData(reader io.ReadSeeker) (HeaderData, error) {
	var headerData HeaderData
	// var buffer []byte = make([]byte, 2)
	var err error

	// Format
	headerData.format, err = parseUint16(reader)

	if err != nil {
		return headerData, err
	}

	// Should be one of 0, 1, 2
	if headerData.format > 2 {
		return headerData, UnsupportedSmfFormat
	}

	// Num tracks
	headerData.numTracks, err = parseUint16(reader)

	if err != nil {
		return headerData, err
	}
	// Division
	var division uint16
	division, err = parseUint16(reader)

	// "If bit 15 of <division> is zero, the bits 14 thru 0 represent the number
	// of delta time "ticks" which make up a quarter-note." 
	if division&0x8000 == 0x0000 {
		headerData.ticksPerQuarterNote = division & 0x7FFF
		headerData.timeFormat = MetricalTimeFormat
	} else {
		// TODO: Can't be bothered to implement this bit just now. 
		// If you want it, write it!
		headerData.timeFormatData = division & 0x7FFF
		headerData.timeFormat = TimeCodeTimeFormat
	}

	if err != nil {
		return headerData, err
	}

	return headerData, nil
}

// readStatusByte reads the track event status byte and returns the type and channel
func readStatusByte(reader io.ReadSeeker) (messageType uint8, messageChannel uint8, err error) {
	var buffer []byte = make([]byte, 1)
	num, err := reader.Read(buffer)

	// If we couldn't read 1 byte, that's a problem.
	if num != 1 {
		return 0, 0, UnexpectedEndOfFile
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return 0, 0, err
	}

	messageType = (buffer[0] & 0xF0) >> 4
	messageChannel = buffer[0] & 0x0F

	return
}

func parseText(reader io.ReadSeeker) (string, error) {
	length, err := parseVarLength(reader)

	if err != nil {
		return "", err
	}

	var buffer []byte = make([]byte, length)

	num, err := reader.Read(buffer)

	// If we couldn't read the entire expected-length buffer, that's a problem.
	if num != int(length) {
		return "", UnexpectedEndOfFile
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return "", err
	}

	// TODO: Data should be ASCII but might go up to 0xFF.
	// What will Go do? Try and decode UTF-8?
	return string(buffer), nil
}