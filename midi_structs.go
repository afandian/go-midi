// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

/*
 * Data structures passed to the Lexer callback.
 */

package midi

// A chunk header
type ChunkHeader struct {
	chunkType string
	length    uint32
}

// Header data
type HeaderData struct {
	format    uint16
	numTracks uint16

	// One of MetricalTimeFormat or TimeCodeTimeFormat
	timeFormat uint

	// Used if TimeCodeTimeFormat
	// Currently data is not un-packed.
	timeFormatData uint16

	// Used if MetricalTimeFormat
	ticksPerQuarterNote uint16
}
