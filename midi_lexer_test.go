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

// Expect a chunk, get MTrk.
// Should store reported track length and go back to ExpectChunk at end of chunk.
// ExpectChunk -> ExpectTrackEvent
func TestMidiLexerShouldHandleChunkLengths(t *testing.T) {
	// TODO
}

// Expect a chunk, get MTrk with a too-short length.
// Should raise a BadSizeChunk error
// ExpectChunk -> ExpectTrackEvent
func TestMidiLexerShouldHandleChunkLengthError(t *testing.T) {
	// TODO
}

// Expect a track event, parse a NoteOff message.
// ExpectTrackEvent -> ExpectTrackEvent
func TestNoteOff(t *testing.T) {
	mockLexerCallback = new(CountingLexerCallback)

	mockReadSeeker = NewMockReadSeeker(&[]byte{0x40, 0x85, 0x04, 0x03})
	lexer = NewMidiLexer(mockReadSeeker, mockLexerCallback)

	// Pre: ExpectChunk
	// Should be ready for a chunk.
	lexer.state = ExpectTrackEvent

	finished, err = lexer.next()
	assertNoError(err, t)

	// Post:
	// not finished yet
	assertFalse(finished, t)

	// ExpectChunk state.
	assertIntsEqual(lexer.state, ExpectTrackEvent, t)

	// callback.Track should have been called.
	assertIntsEqual(mockLexerCallback.noteOff, 1, t)
	assertUint32Equal(mockLexerCallback.time, 0x40, t)
	assertUint8sEqual(mockLexerCallback.channel, 0x05, t)
	assertUint8sEqual(mockLexerCallback.pitch, 0x04, t)
	assertUint8sEqual(mockLexerCallback.velocity, 0x03, t)
}

// Expect a track event, parse a NoteOn message.
// ExpectTrackEvent -> ExpectTrackEvent
func TestNoteOn(t *testing.T) {
	mockLexerCallback = new(CountingLexerCallback)

	mockReadSeeker = NewMockReadSeeker(&[]byte{0x40, 0x95, 0x04, 0x03})
	lexer = NewMidiLexer(mockReadSeeker, mockLexerCallback)

	// Pre: ExpectChunk
	// Should be ready for a chunk.
	lexer.state = ExpectTrackEvent

	finished, err = lexer.next()
	assertNoError(err, t)

	// Post:
	// not finished yet
	assertFalse(finished, t)

	// ExpectChunk state.
	assertIntsEqual(lexer.state, ExpectTrackEvent, t)

	// callback.Track should have been called.
	assertIntsEqual(mockLexerCallback.noteOn, 1, t)
	assertUint32Equal(mockLexerCallback.time, 0x40, t)
	assertUint8sEqual(mockLexerCallback.channel, 0x05, t)
	assertUint8sEqual(mockLexerCallback.pitch, 0x04, t)
	assertUint8sEqual(mockLexerCallback.velocity, 0x03, t)
}

// Expect a track event, parse a NoteOn message.
// ExpectTrackEvent -> ExpectTrackEvent
func TestNotePolyphonicKeyPressure(t *testing.T) {
	mockLexerCallback = new(CountingLexerCallback)

	mockReadSeeker = NewMockReadSeeker(&[]byte{0x40, 0xA7, 0x12, 0x34})
	lexer = NewMidiLexer(mockReadSeeker, mockLexerCallback)

	// Pre: ExpectChunk
	// Should be ready for a chunk.
	lexer.state = ExpectTrackEvent

	finished, err = lexer.next()
	assertNoError(err, t)

	// Post:
	// not finished yet
	assertFalse(finished, t)

	// ExpectChunk state.
	assertIntsEqual(lexer.state, ExpectTrackEvent, t)

	// callback.Track should have been called.
	assertIntsEqual(mockLexerCallback.polyphonicAfterTouch, 1, t)
	assertUint32Equal(mockLexerCallback.time, 0x40, t)
	assertUint8sEqual(mockLexerCallback.channel, 0x07, t)
	assertUint8sEqual(mockLexerCallback.pitch, 0x12, t)
	assertUint8sEqual(mockLexerCallback.pressure, 0x34, t)
}

// Expect a track event, parse a  message.
// ExpectTrackEvent -> ExpectTrackEvent
func TestProgramChange(t *testing.T) {
	// TODO
}

// Expect a track event, parse a  message.
// ExpectTrackEvent -> ExpectTrackEvent
func TestChannelPressure(t *testing.T) {
	mockLexerCallback = new(CountingLexerCallback)

	mockReadSeeker = NewMockReadSeeker(&[]byte{0x40, 0xD8, 0x56})
	lexer = NewMidiLexer(mockReadSeeker, mockLexerCallback)

	// Pre: ExpectChunk
	// Should be ready for a chunk.
	lexer.state = ExpectTrackEvent

	finished, err = lexer.next()
	assertNoError(err, t)

	// Post:
	// not finished yet
	assertFalse(finished, t)

	// ExpectChunk state.
	assertIntsEqual(lexer.state, ExpectTrackEvent, t)

	// callback.Track should have been called.
	assertIntsEqual(mockLexerCallback.channelAfterTouch, 1, t)
	assertUint32Equal(mockLexerCallback.time, 0x40, t)
	assertUint8sEqual(mockLexerCallback.channel, 0x08, t)
	assertUint8sEqual(mockLexerCallback.pressure, 0x56, t)
}

// Expect a track event, parse a  message.
// ExpectTrackEvent -> ExpectTrackEvent
func TestPitchWheel(t *testing.T) {
	mockLexerCallback = new(CountingLexerCallback)
	
	// Three sequential pitch wheel events. NB the value is 14-bit, 
	// split over two bytes, little end first!
	mockReadSeeker = NewMockReadSeeker(&[]byte{
		0x10, 0xE9, 0x00, 0x40, // 0x2000 should be centre
		0x20, 0xE8, 0x34, 0x24, // 0x1234 encoded
		0x50, 0xE7, 0x00, 0x40})
	lexer = NewMidiLexer(mockReadSeeker, mockLexerCallback)

	/* 
	 * FIRST
	 */

	// Pre: ExpectChunk
	// Should be ready for a chunk.
	lexer.state = ExpectTrackEvent

	finished, err = lexer.next()
	assertNoError(err, t)

	// Post:
	// not finished yet
	assertFalse(finished, t)

	// ExpectChunk state.
	assertIntsEqual(lexer.state, ExpectTrackEvent, t)

	// callback.Track should have been called.
	assertIntsEqual(mockLexerCallback.pitchWheel, 1, t)
	assertUint32Equal(mockLexerCallback.time, 0x10, t)
	assertUint8sEqual(mockLexerCallback.channel, 0x09, t)
	assertInt16sEqual(mockLexerCallback.pitchWheelValue, 0x00, t)
	assertUint16Equal(mockLexerCallback.pitchWheelValueAbsolute, 0x2000, t)

	/* 
	 * SECOND
	 */
	 
	// Pre: ExpectChunk
	// Should be ready for a chunk.
	lexer.state = ExpectTrackEvent

	finished, err = lexer.next()
	assertNoError(err, t)

	// Post:
	// not finished yet
	assertFalse(finished, t)

	// ExpectChunk state.
	assertIntsEqual(lexer.state, ExpectTrackEvent, t)

	// callback.Track should have been called.
	assertIntsEqual(mockLexerCallback.pitchWheel, 2, t)
	assertUint32Equal(mockLexerCallback.time, 0x20, t)
	assertUint8sEqual(mockLexerCallback.channel, 0x08, t)
	assertInt16sEqual(mockLexerCallback.pitchWheelValue, -0xDCC, t)
	assertUint16Equal(mockLexerCallback.pitchWheelValueAbsolute, 0x1234, t)

	/* 
	 * THIRD
	 */
	 
	// Pre: ExpectChunk
	// Should be ready for a chunk.
	lexer.state = ExpectTrackEvent

	finished, err = lexer.next()
	assertNoError(err, t)

	// Post:
	// not finished yet
	assertFalse(finished, t)

	// ExpectChunk state.
	assertIntsEqual(lexer.state, ExpectTrackEvent, t)

	// callback.Track should have been called.
	assertIntsEqual(mockLexerCallback.pitchWheel, 3, t)
	assertUint32Equal(mockLexerCallback.time, 0x50, t)
	assertUint8sEqual(mockLexerCallback.channel, 0x07, t)
	assertInt16sEqual(mockLexerCallback.pitchWheelValue, 0x00, t)
	assertUint16Equal(mockLexerCallback.pitchWheelValueAbsolute, 0x2000, t)
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
