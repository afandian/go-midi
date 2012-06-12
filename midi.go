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
	"fmt"
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

	fmt.Println("Lexer next state", lexer.state)

	// See comments for state values above.
	switch lexer.state {
	case ExpectHeader:
		{
			var chunkHeader ChunkHeader
			chunkHeader, err = parseChunkHeader(lexer.input)
			if chunkHeader.chunkType != "MThd" {
				err = ExpectedMthd

				fmt.Println("ChunkHeader error ", err)
				return
			}

			var header, err = parseHeaderData(lexer.input)

			if err != nil {
				fmt.Println("HeaderData error ", err)
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
				fmt.Println("Chunk header error ", err)
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

	case ExpectTrackEvent:
		{
			// If we're at the end of the chunk, change the state.
			if lexer.nextChunkHeader != 0 {

				// The chunk should end exactly on the chunk boundary, really.
				if currentPosition == lexer.nextChunkHeader {
					lexer.state = ExpectChunk
					return false, nil
				} else if currentPosition > lexer.nextChunkHeader {
					fmt.Println("Chunk end error ", err)
					return false, BadSizeChunk
				}
			}

			// Time Delta
			time, err := parseVarLength(lexer.input)
			if err != nil {
				fmt.Println("Time delta error ", err)
				return false, err
			}

			// Message type, Message Channel
			mType, channel, err := readStatusByte(lexer.input)

			switch mType {
			// NoteOff
			case 0x8:
				{
					pitch, velocity, err := parseTwoUint7(lexer.input)

					if err != nil {
						fmt.Println("NoteOff error ", err)
						return false, err
					}

					lexer.callback.NoteOff(channel, pitch, velocity, time)
				}

			// NoteOn
			case 0x9:
				{
					pitch, velocity, err := parseTwoUint7(lexer.input)

					if err != nil {
						fmt.Println("NoteOn error ", err)
						return false, err
					}

					lexer.callback.NoteOn(channel, pitch, velocity, time)
				}

			// Polyphonic Key Pressure
			case 0xA:
				{
					pitch, pressure, err := parseTwoUint7(lexer.input)

					if err != nil {
						return false, err
					}

					lexer.callback.PolyphonicAfterTouch(channel, pitch, pressure, time)
				}

			// Control Change
			case 0xB:
				{
					// channel, value, err := parseTwoUint7(lexer.input)

				}

			// Program Change
			// case 0xC : {
			// 	program, err := parseUint7(lexer.input)

			// 	if err != nil {
			// 		return false, err
			// 		}

			// }

			// Channel Pressure
			case 0xD:
				{
					value, err := parseUint7(lexer.input)

					if err != nil {
						return false, err
					}

					lexer.callback.ChannelAfterTouch(channel, value, time)
				}

			// Pitch Wheel
			case 0xE:
				{
					// The value is a signed int (relative to centre), and absoluteValue is the actual value in the file.
					value, absoluteValue, err := parsePitchWheelValue(lexer.input)

					if err != nil {
						return false, err
					}

					lexer.callback.PitchWheel(channel, value, absoluteValue, time)
				}

				// System Common and System Real-Time
				// case 0xF:
				// 	{
				//		controller, value, err := parseTwoUint7(lexer.input)
				// 	}

				// This covers all cases.
			}

			// Now we need to see if we're at the end of a Track Data chunk.

		}
	}

	return
}
