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

	// Midi messages
	NoteOff(channel uint8, pitch uint8, velocity uint8)
	NoteOn(channel uint8, pitch uint8, velocity uint8)
	PolyphonicAfterTouch(channel uint8, pitch uint8, pressure uint8)
	ControlChange(channel uint8, controller uint8, value uint8)
	ProgramChange(channel uint8, program uint8)
	ChannelAfterTouch(channel uint8, value uint8)
	PitchWheel(channel uint8, value uint16)
	TimeCodeQuarter(messageType uint8, values uint8)
	SongPositionPointer(beats uint16)
	SongSelect(song uint8)
	Undefined1()
	Undefined2()
	TuneRequest()
	TimingClock()
	Undefined3()
	Start()
	Continue()
	Stop()
	Undefined4()
	ActiveSensing()
	Reset()
	Done()
}
