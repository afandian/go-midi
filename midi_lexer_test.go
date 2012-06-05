// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

// Tests for lexer. Slightly higher level.

package midi

import (
	"io"
	"testing"
)

var lexer *MidiLexer
var mockLexerCallback *CountingLexerCallback
var mockReadSeeker io.ReadSeeker
var finished bool
var err error

// Get clean values.
func setupData(data *[]byte) {
	mockLexerCallback = new(CountingLexerCallback)
	mockReadSeeker = NewMockReadSeeker(data)
}

// MidiLexer should throw error for null callback or input
func TestLexerShouldComplainNullArgs(t *testing.T) {
	setupData(&[]byte{})

	var status int

	// First call with good arguments.
	lexer = NewMidiLexer(mockReadSeeker, mockLexerCallback)
	status = lexer.Lex()
	if status != Ok {
		t.Fatal("Status should be OK")
	}

	// Call with no reader
	lexer = NewMidiLexer(nil, mockLexerCallback)
	status = lexer.Lex()
	assertHasFlag(status, NoReadSeeker, t)

	// Call with no callback
	lexer = NewMidiLexer(mockReadSeeker, nil)
	status = lexer.Lex()
	assertHasFlag(status, NoCallback, t)
}

/*
 * Correct state transitions. 
 */

// Start of file, consume header.
// ExpectHeader -> ExpectChunk
func TestLexerShouldExpectHeader(t *testing.T) {
	// Just enough for the header chunk

	mockLexerCallback = new(CountingLexerCallback)

	mockReadSeeker = NewMockReadSeeker(&[]byte{0x4D, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x00, 0xC8})

	lexer = NewMidiLexer(mockReadSeeker, mockLexerCallback)

	// Pre: New file, ExpectHeader state.
	// Should be ready for header
	assertIntsEqual(lexer.state, ExpectHeader, t)

	finished, err = lexer.next()
	assertNoError(err, t)

	// Post:
	// not finished yet
	assertFalse(finished, t)

	// ExpectChunk state.
	assertIntsEqual(lexer.state, ExpectChunk, t)

	// Began() was called.
	assertIntsEqual(mockLexerCallback.began, 1, t)

	// Began() was called with the right values.
	assertIntsEqual(int(mockLexerCallback.header), 1, t)
	assertIntsEqual(int(mockLexerCallback.headerData.format), 1, t)
	assertIntsEqual(int(mockLexerCallback.headerData.numTracks), 2, t)

	if mockLexerCallback.headerData.timeFormat != MetricalTimeFormat {
		t.Fatal("Was not MetricalTimeFormat")
	}

	assertIntsEqual(int(mockLexerCallback.headerData.ticksPerQuarterNote), 200, t)
}

// Expect a chunk, get an unrecognised type. Should skip to next.
// ExpectChunk -> ExpectChunk
func TestMidiLexerShouldSkipUnknownTrack(t *testing.T) {
	// Just enough for the header chunk

	mockLexerCallback = new(CountingLexerCallback)

	// Head of data stream is MThd, where the lexer will expect MTrk

	mockReadSeeker = NewMockReadSeeker(&[]byte{ /* start of unknown block, claims to be 2-long */ 0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x00, 0x00, 0x02, 0xCA, 0xFE /* Start of next block. */, 0x4D, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x00, 0xC8})

	lexer = NewMidiLexer(mockReadSeeker, mockLexerCallback)

	// Pre: ExpectChunk
	// Should be ready for a chunk.
	lexer.state = ExpectChunk

	finished, err = lexer.next()
	assertNoError(err, t)

	// Post:
	// not finished yet
	assertFalse(finished, t)

	// ExpectChunk state.
	assertIntsEqual(lexer.state, ExpectChunk, t)

	// Reader should have jumped to position 10, the next block.
	var position, err = lexer.input.Seek(0, 1)
	assertNoError(err, t)
	assertIntsEqual(int(position), 10, t)
}

// Expect a chunk, get MTrk. Should enter ExpectTrackEvent state.
// ExpectChunk -> ExpectTrackEvent
func TestMidiLexerShouldExpectTrackEvent(t *testing.T) {
	mockLexerCallback = new(CountingLexerCallback)

	// Head of data stream is MThd, where the lexer will expect MTrk

	mockReadSeeker = NewMockReadSeeker(&[]byte{0x4D, 0x54, 0x72, 0x6B, 0x00, 0x00, 0x00, 0xEE, 0x00, 0x01, 0x00, 0x02, 0x00, 0xC8})

	lexer = NewMidiLexer(mockReadSeeker, mockLexerCallback)

	// Pre: ExpectChunk
	// Should be ready for a chunk.
	lexer.state = ExpectChunk

	finished, err = lexer.next()
	assertNoError(err, t)

	// Post:
	// not finished yet
	assertFalse(finished, t)

	// ExpectChunk state.
	assertIntsEqual(lexer.state, ExpectTrackEvent, t)

	// callback.Track should have been called.
	assertIntsEqual(mockLexerCallback.track, 1, t)
	assertIntsEqual(int(mockLexerCallback.chunkHeader.length), 0xEE, t)
}

/*
 * Exceptional state transitions. 
 */

// Bad header chunk type at start of file should result in error
func TestLexerShouldErrorBadHeader(t *testing.T) {

	// Just enough for the header chunk
	mockLexerCallback = new(CountingLexerCallback)

	mockReadSeeker = NewMockReadSeeker(&[]byte{0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x00, 0xC8})

	lexer = NewMidiLexer(mockReadSeeker, mockLexerCallback)

	finished, err = lexer.next()

	assertError(err, ExpectedMthd, t)
}
