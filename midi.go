// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

/*
 * The main file.
 * This contains the lexer, which does the job of scanning through the file.
 */

package midi

import (
	"io"
)

// State of the MidiLexerCallback.
const (
	// At the start of the MIDI file.
	// Expect SMF Header chunk.
	ExpectHeader = iota

	// Expect a chunk. Any kind of chunk. Except MThd.
	// But really, anything other than MTrk would be weird.
	ExpectChunk = iota

	// We're in a Track, expect a track event.
	ExpectTrackEvent = iota

	// This has to happen sooner or later.
	Done = iota
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
	return &MidiLexer{callback: callback, input: input, state: ExpectHeader}
}

// Lex starts the MidiLexer running.
func (lexer *MidiLexer) Lex() error {
	if lexer.callback == nil {
		return NoCallback
	}

	if lexer.input == nil {
		return NoReadSeeker
	}

	var finished bool = false
	var err error

	for {
		finished, err = lexer.next()

		if err != nil {
			return err
		}

		if finished == true {
			return nil
		}
	}

	return nil
}

// next lexes the next item, calling appropriate callbacks.
// Finished only set true when finished correctly.
func (lexer *MidiLexer) next() (finished bool, err error) {

	// Default for return values.
	err = nil
	finished = false

	// The position in the file before the next lexing event happens.
	// Useful in some cases
	var currentPosition int64
	currentPosition, err = lexer.input.Seek(0, 1)

	if err != nil {
		return
	}

	// See comments for state values above.
	switch lexer.state {
	case ExpectHeader:
		{
			//fmt.Println("ExpectHeader")

			var chunkHeader ChunkHeader
			chunkHeader, err = parseChunkHeader(lexer.input)
			if chunkHeader.ChunkType != "MThd" {
				err = ExpectedMthd

				//fmt.Println("ChunkHeader error ", err)
				return
			}

			var header HeaderData
			header, err = parseHeaderData(lexer.input)

			if err != nil {
				//fmt.Println("HeaderData error ", err)
				return
			}

			lexer.callback.Began()

			lexer.callback.Header(header)

			lexer.state = ExpectChunk

			return
		}

	case ExpectChunk:
		{
			//fmt.Println("ExpectChunk")

			var chunkHeader ChunkHeader
			chunkHeader, err = parseChunkHeader(lexer.input)

			//fmt.Println("Got chunk header", chunkHeader)

			if err != nil {
				// If we expect a chunk and we hit the end of the file, that's not so unexpected after all.
				// The file has to end some time, and this is the correct boundary upon which to end it.
				if err == UnexpectedEndOfFile {
					lexer.state = Done

					// TODO TEST
					lexer.callback.Finished()

					finished = true
					err = nil
					return
				}

				//fmt.Println("Chunk header error ", err)
				return
			}

			lexer.callback.Track(chunkHeader)
			lexer.nextChunkHeader = int64(chunkHeader.Length) + currentPosition

			// If the header is of an unknown type, skip over it.
			if chunkHeader.ChunkType != "MTrk" {
				lexer.input.Seek(lexer.nextChunkHeader, 1)

				// Then we expect another chunk.
				lexer.state = ExpectChunk
				lexer.nextChunkHeader = 0
			} else {
				// We have a MTrk
				lexer.state = ExpectTrackEvent
			}

			return
		}

	case ExpectTrackEvent:
		{
			//fmt.Println("ExpectTrackEvent")

			// Removed because there is an event to say 'end of chunk'.
			// TODO: investigate. Could put this back for error cases.
			// // If we're at the end of the chunk, change the state.
			// if lexer.nextChunkHeader != 0 {

			// 	// The chunk should end exactly on the chunk boundary, really.
			// 	if currentPosition == lexer.nextChunkHeader {
			// 		lexer.state = ExpectChunk
			// 		return false, nil
			// 	} else if currentPosition > lexer.nextChunkHeader {
			// 		//fmt.Println("Chunk end error ", err)
			// 		return false, BadSizeChunk
			// 	}
			// }

			// Time Delta
			var time uint32
			time, err = parseVarLength(lexer.input)
			if err != nil {
				//fmt.Println("Time delta error ", err)
				return
			}

			// Message type, Message Channel
			var mType, channel uint8
			mType, channel, err = readStatusByte(lexer.input)

			//fmt.Println("Track Event Type ", mType)

			switch mType {
			// NoteOff
			case 0x8:
				{
					var pitch, velocity uint8
					pitch, velocity, err = parseTwoUint7(lexer.input)

					if err != nil {
						//fmt.Println("NoteOff error ", err)
						return
					}

					lexer.callback.NoteOff(channel, pitch, velocity, time)
				}

			// NoteOn
			case 0x9:
				{
					var pitch, velocity uint8
					pitch, velocity, err = parseTwoUint7(lexer.input)

					if err != nil {
						//fmt.Println("NoteOn error ", err)
						return
					}

					lexer.callback.NoteOn(channel, pitch, velocity, time)
				}

			// Polyphonic Key Pressure
			case 0xA:
				{
					var pitch, pressure uint8
					pitch, pressure, err = parseTwoUint7(lexer.input)

					if err != nil {
						return
					}

					lexer.callback.PolyphonicAfterTouch(channel, pitch, pressure, time)
				}

			// Control Change or Channel Mode Message
			case 0xB:
				{
					var controller, value uint8
					controller, value, err = parseTwoUint7(lexer.input)

					if err != nil {
						return
					}

					// TODO split this into ChannelMode for values [120, 127]?
					// TODO implement separate callbacks for each type of:
					// - All sound off
					// - Reset all controllers
					// - Local control
					// - All notes off
					// Only if required. http://www.midi.org/techspecs/midimessages.php

					// TODO TEST
					lexer.callback.ControlChange(channel, controller, value, time)
					return
				}

			// Program Change
			case 0xC:
				{
					var program uint8
					program, err = parseUint7(lexer.input)

					if err != nil {
						return
					}

					lexer.callback.ProgramChange(channel, program, time)
					return
				}

			// Channel Pressure
			case 0xD:
				{
					var value uint8
					value, err = parseUint7(lexer.input)

					if err != nil {
						return
					}

					lexer.callback.ChannelAfterTouch(channel, value, time)
				}

			// Pitch Wheel
			case 0xE:
				{
					// The value is a signed int (relative to centre), and absoluteValue is the actual value in the file.
					var value int16
					var absoluteValue uint16
					value, absoluteValue, err = parsePitchWheelValue(lexer.input)

					if err != nil {
						return
					}

					lexer.callback.PitchWheel(channel, value, absoluteValue, time)
				}

			// System Common and System Real-Time / Meta
			case 0xF:
				{
					// The 4-bit nibble called 'channel' isn't actually the channel in this case.
					switch channel {
					// Meta-events
					case 0xF:
						{
							var command uint8
							command, err = parseUint8(lexer.input)

							if err != nil {
								return
							}

							//fmt.Println("SystemCommon/RealTime command:", command)

							// TODO: If every one of these takes a length, then take this outside.
							// Will make for more more robust unknown types.
							switch command {

							// Sequence number
							case 0x00:
								{
									//fmt.Println("seq no")
									var length uint8
									length, err = parseUint8(lexer.input)

									if err != nil {
										return
									}

									// Zero length sequences allowed according to http://home.roadrunner.com/~jgglatt/tech/midifile/seq.htm
									if length == 0 {
										lexer.callback.SequenceNumber(channel, 0, false, time)

										return
									}

									// Otherwise length will be 2 to hold the uint16.
									var sequenceNumber uint16
									sequenceNumber, err = parseUint16(lexer.input)

									if err != nil {
										return
									}

									lexer.callback.SequenceNumber(channel, sequenceNumber, true, time)

									return
								}

							// Text event
							case 0x01:
								{
									//fmt.Println("Text")
									var text string
									text, err = parseText(lexer.input)
									//fmt.Println("text value", text, err)
									if err != nil {
										return
									}

									lexer.callback.Text(channel, text, time)

									return
								}

							// Copyright text event
							case 0x02:
								{
									//fmt.Println("Copyright")
									var text string
									text, err = parseText(lexer.input)

									if err != nil {
										return
									}

									lexer.callback.CopyrightText(channel, text, time)

									return
								}

							// Sequence or track name
							case 0x03:
								{
									var text string
									text, err = parseText(lexer.input)

									if err != nil {
										return
									}

									lexer.callback.SequenceName(channel, text, time)

									return

								}

							// Track instrument name
							case 0x04:
								{
									var text string
									text, err = parseText(lexer.input)

									if err != nil {
										return
									}

									lexer.callback.TrackInstrumentName(channel, text, time)

									return

								}

							// Lyric text
							case 0x05:
								{
									var text string
									text, err = parseText(lexer.input)

									if err != nil {
										return
									}

									lexer.callback.LyricText(channel, text, time)

									return
								}

							// Marker text
							case 0x06:
								{
									var text string
									text, err = parseText(lexer.input)

									if err != nil {
										return
									}

									lexer.callback.MarkerText(channel, text, time)

									return
								}

							// Cue point text
							case 0x07:
								{
									var text string
									text, err = parseText(lexer.input)

									if err != nil {
										return
									}

									lexer.callback.CuePointText(channel, text, time)

									return
								}

							case 0x20:
								{
									// Obsolete 'MIDI Channel'
									//fmt.Println("MIDI Channel obsolete")
									var length uint32
									length, err = parseVarLength(lexer.input)

									if err != nil {
										return
									}

									if length != 1 {
										err = UnexpectedEventLengthError{"Midi Channel Event expected length 1"}
										return
									}

									// This is the channel value.
									// Just forget this one.
									_, err = parseUint8(lexer.input)

									if err != nil {
										return
									}
								}

							case 0x21:
								{
									// Obsolete 'MIDI Port'
									//fmt.Println("MIDI PORT obsolete")
									var length uint32
									length, err = parseVarLength(lexer.input)

									if err != nil {
										return
									}

									if length != 1 {
										err = UnexpectedEventLengthError{"MIDI Port Event expected length 1"}
										return
									}

									// This is the port value.
									// Just forget this one.
									_, err = parseUint8(lexer.input)

									if err != nil {
										return
									}
								}

							// End of track
							case 0x2F:
								{
									var length uint32
									length, err = parseVarLength(lexer.input)

									if err != nil {
										return
									}

									if length != 0 {
										err = UnexpectedEventLengthError{"EndOfTrack expected length 0"}
										return
									}

									lexer.callback.EndOfTrack(channel, time)

									// Expect the next chunk event.
									lexer.state = ExpectChunk

									return false, nil
								}

							// Set tempo
							case 0x51:
								{
									// TODO TEST

									var length uint32
									length, err = parseVarLength(lexer.input)

									if err != nil {
										return
									}

									if length != 3 {
										err = UnexpectedEventLengthError{"Tempo expected length 3"}
										return
									}

									var microsecondsPerCrotchet uint32
									microsecondsPerCrotchet, err = parseUint24(lexer.input)

									if err != nil {
										return
									}

									// Also beats per minute
									var bpm uint32
									bpm = 60000000 / microsecondsPerCrotchet

									lexer.callback.Tempo(bpm, microsecondsPerCrotchet, time)
								}

							// Time signature
							case 0x58:
								{
									var length uint32
									length, err = parseVarLength(lexer.input)

									if err != nil {
										return
									}

									if length != 4 {
										err = UnexpectedEventLengthError{"TimeSignature expected length 4"}
										return
									}

									// TODO TEST
									var numerator uint8
									numerator, err = parseUint8(lexer.input)

									if err != nil {
										return
									}

									var denomenator uint8
									denomenator, err = parseUint8(lexer.input)

									if err != nil {
										return
									}

									var clocksPerClick uint8
									clocksPerClick, err = parseUint8(lexer.input)

									if err != nil {
										return
									}

									var demiSemiQuaverPerQuarter uint8
									demiSemiQuaverPerQuarter, err = parseUint8(lexer.input)

									if err != nil {
										return
									}

									//fmt.Println("TimeSignature event", numerator, denomenator, clocksPerClick, demiSemiQuaverPerQuarter, time)

									lexer.callback.TimeSignature(numerator, denomenator, clocksPerClick, demiSemiQuaverPerQuarter, time)

									return false, nil
								}

							// Key signature
							case 0x59:
								{
									// TODO TEST
									var length uint32
									var sharpsOrFlats int8
									var mode uint8

									length, err = parseVarLength(lexer.input)

									if err != nil {
										return
									}

									if length != 2 {
										err = UnexpectedEventLengthError{"KeySignature expected length 2"}
										return
									}

									// Signed int, positive is sharps, negative is flats.
									sharpsOrFlats, err = parseInt8(lexer.input)

									if err != nil {
										return
									}

									// Mode is Major or Minor.
									mode, err = parseUint8(lexer.input)

									if err != nil {
										return
									}

									key, resultMode := keySignatureFromSharpsOrFlats(sharpsOrFlats, mode)

									lexer.callback.KeySignature(key, resultMode, sharpsOrFlats)
								}

							// Sequencer specific info
							case 0x7F:
								{
									//fmt.Println("0x7F")
								}

							// Timing clock
							case 0xF8:
								{
									//fmt.Println("0xF8")
								}

							// Start current sequence
							case 0xFA:
								{
									//fmt.Println("0xFA")
								}

							// Continue stopped sequence where left off
							case 0xFB:
								{
									//fmt.Println("0xFB")
								}

							// Stop sequence
							case 0xFc:
								{
									//fmt.Println("0xFc")
								}
							default:
								//fmt.Println("Unrecognised meta command", command)
							}

						}

					default:
						//fmt.Println("Unrecognised message type", mType)
					}

					//
				}

				// This covers all cases.

			// Now we need to see if we're at the end of a Track Data chunk.
			default:
				{
					var length uint32
					length, err = parseVarLength(lexer.input)

					if err != nil {
						return
					}

					//fmt.Println("Type Unrecognised", mType, "length", length)

					// Read length of chunk
					for i := uint32(0); i < length; i++ {
						_, err = parseUint8(lexer.input)

						if err != nil {
							return
						}
					}
				}
			}

		}

	case Done:
		{
			// The event that raised this will already have returned false to say it's stopped ticking.
			// Just keep returning false.
			finished = true
			return
		}
	}

	return
}
