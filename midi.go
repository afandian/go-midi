// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

// Main file

package midi

import (
	"io"
)

// State of the MidiLexer.
const (
	// At the start of the MIDI file.
	// Expect SMF Header chunk.
	ExpectHeader = iota

	// Expect a chunk. Any kind of chunk. Except MThd.
	// But really, anything other than MTrk would be weird.
	ExpectChunk = iota

	// We're in a Track, expect a track event.
	ExpectTrackEvent = iota
)

// MidiLexer is a Standard Midi File Lexer.
// Pass this a ReadSeeker to a MIDI file and a callback that conforms to MidiLexerCallback 
// and it'll run over the file, calling events on the callback.
type MidiLexer struct {
	callback MidiLexerCallback
	input    io.ReadSeeker

	// State of the parser, as per the above constants.
	state int

	// The location of the next chunk header that we expect to find as an offset from
	// most recent Chunk header.
	nextChunkHeader int64
}

// Construct a new MidiLexer
func NewMidiLexer(input io.ReadSeeker, callback MidiLexerCallback) *MidiLexer {
	return &MidiLexer{callback: callback, input: input}
}

// Lex starts the MidiLexer running.
func (lexer *MidiLexer) Lex() (error int) {
	if lexer.callback == nil {
		return NoCallback
	}

	if lexer.input == nil {
		return NoReadSeeker
	}

	return Ok
}

// next lexes the next item, calling appropriate callbacks.
// Finished only set true when finished correctly.
func (lexer *MidiLexer) next() (finished bool, err error) {
	err = nil
	finished = false

	// The position in the file before the next lexing event happens. 
	// Useful in some cases
	currentPosition, err := lexer.input.Seek(0, 1)

	if err != nil {
		return
	}

	// See comments for state values above.
	switch lexer.state {
	case ExpectHeader:
		{
			var chunkHeader ChunkHeader
			chunkHeader, err = parseChunkHeader(lexer.input)
			if chunkHeader.chunkType != "MThd" {
				err = ExpectedMthd
				return
			}

			var header, err = parseHeaderData(lexer.input)

			if err != nil {
				return false, err
			}

			lexer.callback.Began()
			lexer.callback.Header(header)

			lexer.state = ExpectChunk
		}

	case ExpectChunk:
		{
			var chunkHeader, err = parseChunkHeader(lexer.input)

			if err != nil {
				return false, err
			}

			lexer.callback.Track(chunkHeader)
			lexer.nextChunkHeader = int64(chunkHeader.length) + currentPosition

			// If the header is of an unknown type, skip over it.
			if chunkHeader.chunkType != "MTrk" {
				lexer.input.Seek(lexer.nextChunkHeader, 1)

				// Then we expect another chunk.
				lexer.state = ExpectChunk
				lexer.nextChunkHeader = 0
			} else {
				// We have a MTrk
				lexer.state = ExpectTrackEvent
			}
		}

		// case ExpectTrackEvent: {
		// 	timeDelta, err := parseVarLength(lexer.input)
		// }
	}

	return
}
