// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

// Interface

package midi

// MidiLexerCallback describes a callback object that should be passed to MidiLexer. It recieves the following method calls as the lexer finds them.
type MidiLexerCallback interface {
	// Meta messages

	// Started reading a file.
	Began()

	// Finished reading the file.
	Finished()

	// There was an error when lexing.
	ErrorReading()

	// There was an error opening the file input.
	ErrorOpeningFile()

	// SMF header.
	Header(header HeaderData)

	// A chunk header (usually MTrk).
	Track(header ChunkHeader)

	// Midi in-track messages
	NoteOff(channel uint8, pitch uint8, velocity uint8, time uint32)
	NoteOn(channel uint8, pitch uint8, velocity uint8, time uint32)
	PolyphonicAfterTouch(channel uint8, pitch uint8, pressure uint8, time uint32)
	ControlChange(channel uint8, controller uint8, value uint8, time uint32)
	ProgramChange(channel uint8, program uint8, time uint32)
	ChannelAfterTouch(channel uint8, value uint8, time uint32)
	PitchWheel(channel uint8, value int16, absValue uint16, time uint32)
	TimeCodeQuarter(messageType uint8, values uint8, time uint32)
	SongPositionPointer(beats uint16, time uint32)
	SongSelect(song uint8, time uint32)
	Undefined1(time uint32)
	Undefined2(time uint32)
	TuneRequest(time uint32)
	TimingClock(time uint32)
	Undefined3(time uint32)
	Start(time uint32)
	Continue(time uint32)
	Stop(time uint32)
	Undefined4(time uint32)
	ActiveSensing(time uint32)
	Reset(time uint32)
	Done(time uint32)
}
