// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

/*
 * Mocks implementations for testing.
 * In order to check that the right calls are made to the LexerCallback by the Lexer, a mock callback is pssed to it during tests.
 */

package midi

import (
	"io"
)

// A mock implementation of LexerCallback that does nothing.
type MockLexerCallback struct{}

func (*MockLexerCallback) Header(header HeaderData)                                        {}
func (*MockLexerCallback) Track(header ChunkHeader)                                        {}
func (*MockLexerCallback) Began()                                                          {}
func (*MockLexerCallback) Finished()                                                       {}
func (*MockLexerCallback) ErrorReading()                                                   {}
func (*MockLexerCallback) ErrorOpeningFile()                                               {}
func (*MockLexerCallback) NoteOff(channel uint8, pitch uint8, velocity uint8, time uint32) {}
func (*MockLexerCallback) NoteOn(channel uint8, pitch uint8, velocity uint8, time uint32)  {}
func (*MockLexerCallback) PolyphonicAfterTouch(channel uint8, pitch uint8, pressure uint8, time uint32) {
}
func (*MockLexerCallback) ControlChange(channel uint8, controller uint8, value uint8, time uint32) {}
func (*MockLexerCallback) ProgramChange(channel uint8, program uint8, time uint32)                 {}
func (*MockLexerCallback) ChannelAfterTouch(channel uint8, value uint8, time uint32)               {}
func (*MockLexerCallback) PitchWheel(channel uint8, value int16, absValue uint16, time uint32)     {}
func (*MockLexerCallback) TimeCodeQuarter(messageType uint8, values uint8, time uint32)            {}
func (*MockLexerCallback) SongPositionPointer(beats uint16, time uint32)                           {}
func (*MockLexerCallback) SongSelect(song uint8, time uint32)                                      {}
func (*MockLexerCallback) Undefined1(time uint32)                                                  {}
func (*MockLexerCallback) Undefined2(time uint32)                                                  {}
func (*MockLexerCallback) TuneRequest(time uint32)                                                 {}
func (*MockLexerCallback) TimingClock(time uint32)                                                 {}
func (*MockLexerCallback) Undefined3(time uint32)                                                  {}
func (*MockLexerCallback) Start(time uint32)                                                       {}
func (*MockLexerCallback) Continue(time uint32)                                                    {}
func (*MockLexerCallback) Stop(time uint32)                                                        {}
func (*MockLexerCallback) Undefined4(time uint32)                                                  {}
func (*MockLexerCallback) ActiveSensing(time uint32)                                               {}
func (*MockLexerCallback) Reset(time uint32)                                                       {}
func (*MockLexerCallback) Done(time uint32)                                                        {}

func (*MockLexerCallback) SequenceNumber(channel uint8, number uint16, numberGiven bool, time uint32) {
}
func (*MockLexerCallback) Text(channel uint8, text string, time uint32)                {}
func (*MockLexerCallback) CopyrightText(channel uint8, text string, time uint32)       {}
func (*MockLexerCallback) SequenceName(channel uint8, text string, time uint32)        {}
func (*MockLexerCallback) TrackInstrumentName(channel uint8, text string, time uint32) {}
func (*MockLexerCallback) LyricText(channel uint8, text string, time uint32)           {}
func (*MockLexerCallback) MarkerText(channel uint8, text string, time uint32)          {}
func (*MockLexerCallback) CuePointText(channel uint8, text string, time uint32)        {}
func (*MockLexerCallback) EndOfTrack(channel uint8, time uint32)                       {}

// A mock implementation of LexerCallback that counts each method call and stores the most recent values,
// so that calls can be verified.
type CountingLexerCallback struct {

	// Callback counts
	header               int
	track                int
	began                int
	finished             int
	errorReading         int
	errorOpeningFile     int
	noteOff              int
	noteOn               int
	polyphonicAfterTouch int
	controlChange        int
	programChange        int
	channelAfterTouch    int
	pitchWheel           int
	timeCodeQuarter      int
	songPositionPointer  int
	songSelect           int
	undefined1           int
	undefined2           int
	tuneRequest          int
	timingClock          int
	undefined3           int
	start                int
	continue_            int
	stop                 int
	undefined4           int
	activeSensing        int
	reset                int
	done                 int
	endOfTrack           int
	text                 int
	copyrightText        int
	sequenceName         int
	trackInstrumentName  int
	lyricText            int
	markerText           int
	cuePointText         int
	sequenceNumber       int

	// Most recent values
	headerData  HeaderData
	chunkHeader ChunkHeader
	pitch       uint8
	channel     uint8
	time        uint32
	velocity    uint8
	pressure    uint8
	textValue   string

	pitchWheelValue         int16
	pitchWheelValueAbsolute uint16
	sequenceNumberGiven     bool
	sequenceNumberValue     uint16
}

func (cbk *CountingLexerCallback) Header(header HeaderData) { cbk.header++; cbk.headerData = header }
func (cbk *CountingLexerCallback) Track(header ChunkHeader) { cbk.track++; cbk.chunkHeader = header }
func (cbk *CountingLexerCallback) Began()                   { cbk.began++ }
func (cbk *CountingLexerCallback) Finished()                { cbk.finished++ }
func (cbk *CountingLexerCallback) ErrorReading()            { cbk.errorReading++ }
func (cbk *CountingLexerCallback) ErrorOpeningFile()        { cbk.errorOpeningFile++ }
func (cbk *CountingLexerCallback) NoteOff(channel uint8, pitch uint8, velocity uint8, time uint32) {
	cbk.noteOff++
	cbk.pitch = pitch
	cbk.channel = channel
	cbk.velocity = velocity
	cbk.time = time
}
func (cbk *CountingLexerCallback) NoteOn(channel uint8, pitch uint8, velocity uint8, time uint32) {
	cbk.noteOn++
	cbk.pitch = pitch
	cbk.channel = channel
	cbk.velocity = velocity
	cbk.time = time
}
func (cbk *CountingLexerCallback) PolyphonicAfterTouch(channel uint8, pitch uint8, pressure uint8, time uint32) {
	cbk.polyphonicAfterTouch++
	cbk.channel = channel
	cbk.pitch = pitch
	cbk.pressure = pressure
	cbk.time = time
}
func (cbk *CountingLexerCallback) ControlChange(channel uint8, controller uint8, value uint8, time uint32) {
	cbk.controlChange++
}
func (cbk *CountingLexerCallback) ProgramChange(channel uint8, program uint8, time uint32) {
	cbk.programChange++
}
func (cbk *CountingLexerCallback) ChannelAfterTouch(channel uint8, pressure uint8, time uint32) {
	cbk.channelAfterTouch++
	cbk.channel = channel
	cbk.pressure = pressure
	cbk.time = time
}
func (cbk *CountingLexerCallback) PitchWheel(channel uint8, value int16, absValue uint16, time uint32) {
	cbk.pitchWheel++
	cbk.pitchWheelValue = value
	cbk.pitchWheelValueAbsolute = absValue
	cbk.time = time
	cbk.channel = channel
}
func (cbk *CountingLexerCallback) TimeCodeQuarter(messageType uint8, values uint8, time uint32) {
	cbk.timeCodeQuarter++
}
func (cbk *CountingLexerCallback) SongPositionPointer(beats uint16, time uint32) {
	cbk.songPositionPointer++
}
func (cbk *CountingLexerCallback) SongSelect(song uint8, time uint32) { cbk.songSelect++ }
func (cbk *CountingLexerCallback) Undefined1(time uint32)             { cbk.undefined1++ }
func (cbk *CountingLexerCallback) Undefined2(time uint32)             { cbk.undefined2++ }
func (cbk *CountingLexerCallback) TuneRequest(time uint32)            { cbk.tuneRequest++ }
func (cbk *CountingLexerCallback) TimingClock(time uint32)            { cbk.timingClock++ }
func (cbk *CountingLexerCallback) Undefined3(time uint32)             { cbk.undefined3++ }
func (cbk *CountingLexerCallback) Start(time uint32)                  { cbk.start++ }
func (cbk *CountingLexerCallback) Continue(time uint32)               { cbk.continue_++ }
func (cbk *CountingLexerCallback) Stop(time uint32)                   { cbk.stop++ }
func (cbk *CountingLexerCallback) Undefined4(time uint32)             { cbk.undefined4++ }
func (cbk *CountingLexerCallback) ActiveSensing(time uint32)          { cbk.activeSensing++ }
func (cbk *CountingLexerCallback) Reset(time uint32)                  { cbk.reset++ }
func (cbk *CountingLexerCallback) Done(time uint32)                   { cbk.done++ }

func (cbk *CountingLexerCallback) SequenceNumber(channel uint8, number uint16, numberGiven bool, time uint32) {
	cbk.sequenceNumber++
	cbk.time = time
	cbk.sequenceNumberValue = number
	cbk.sequenceNumberGiven = numberGiven
}
func (cbk *CountingLexerCallback) Text(channel uint8, text string, time uint32) {
	cbk.text++
	cbk.textValue = text
	cbk.time = time
}
func (cbk *CountingLexerCallback) CopyrightText(channel uint8, text string, time uint32) {
	cbk.copyrightText++
	cbk.textValue = text
	cbk.time = time
}
func (cbk *CountingLexerCallback) SequenceName(channel uint8, text string, time uint32) {
	cbk.sequenceName++
	cbk.textValue = text
	cbk.time = time
}
func (cbk *CountingLexerCallback) TrackInstrumentName(channel uint8, text string, time uint32) {
	cbk.trackInstrumentName++
	cbk.textValue = text
	cbk.time = time
}
func (cbk *CountingLexerCallback) LyricText(channel uint8, text string, time uint32) {
	cbk.lyricText++
	cbk.textValue = text
	cbk.time = time
}
func (cbk *CountingLexerCallback) MarkerText(channel uint8, text string, time uint32) {
	cbk.markerText++
	cbk.textValue = text
	cbk.time = time
}
func (cbk *CountingLexerCallback) CuePointText(channel uint8, text string, time uint32) {
	cbk.cuePointText++
	cbk.textValue = text
	cbk.time = time
}
func (cbk *CountingLexerCallback) EndOfTrack(channel uint8, time uint32) {
	cbk.endOfTrack++
	cbk.time = time
}

// MockReadSeeker is a mock Reader and Seeker. Constructed with data, behaves as a file reader.
// This is used to pass MIDI data to the Lexer and also to the MIDI value parsing functions.
type MockReadSeeker struct {
	data     *[]byte
	position int64
}

// NewMockReadSeeker creates a new MockReadSeeker object backed by the given byte array data.
func NewMockReadSeeker(data *[]byte) *MockReadSeeker {
	return &MockReadSeeker{data: data}
}

// Read fills the given buffer, returning the number of bytes and an error.
func (reader *MockReadSeeker) Read(p []byte) (n int, err error) {
	var amount = int64(len(p))
	var maxAmount = int64(len(*reader.data)) - reader.position

	// Don't read past the end
	if amount > maxAmount {
		amount = maxAmount
	}

	copy(p, (*reader.data)[reader.position:reader.position+amount])
	reader.position += amount
	return int(amount), nil
}

// Seek sets the offset for the next Read or Write to offset, interpreted according to the value of `whence`: 
// 0 means relative to the origin of the file, 1 means relative to the current offset, and 2 means relative to the end.
// Seek returns the new offset and an Error, if any.
func (reader *MockReadSeeker) Seek(offset int64, whence int) (ret int64, err error) {
	switch whence {
	case 0:
		{
			if offset > int64(len(*reader.data)) {
				return -1, io.EOF
			}

			reader.position = offset

			return reader.position, nil
		}
	case 1:
		{
			if offset+reader.position > int64(len(*reader.data)) {
				return -1, io.EOF
			}

			reader.position += offset

			return reader.position, nil
		}
	case 2:
		{
			if offset > int64(len(*reader.data)) {
				return -1, io.EOF
			}

			reader.position = int64(len(*reader.data)) - offset

			return reader.position, nil
		}
	}

	return
}
