// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

// Mocks implementations for testing.

package midi

import (
	"io"
)

// A mock implementation of LexerCallback that does nothing
type MockLexerCallback struct{}

func (*MockLexerCallback) Header(header HeaderData)                                        {}
func (*MockLexerCallback) Track(header ChunkHeader)                                        {}
func (*MockLexerCallback) Began()                                                          {}
func (*MockLexerCallback) Finished()                                                       {}
func (*MockLexerCallback) ErrorReading()                                                   {}
func (*MockLexerCallback) ErrorOpeningFile()                                               {}
func (*MockLexerCallback) NoteOff(channel uint8, pitch uint8, velocity uint8)              {}
func (*MockLexerCallback) NoteOn(channel uint8, pitch uint8, velocity uint8)               {}
func (*MockLexerCallback) PolyphonicAfterTouch(channel uint8, pitch uint8, pressure uint8) {}
func (*MockLexerCallback) ControlChange(channel uint8, controller uint8, value uint8)      {}
func (*MockLexerCallback) ProgramChange(channel uint8, program uint8)                      {}
func (*MockLexerCallback) ChannelAfterTouch(channel uint8, value uint8)                    {}
func (*MockLexerCallback) PitchWheel(channel uint8, value uint16)                          {}
func (*MockLexerCallback) TimeCodeQuarter(messageType uint8, values uint8)                 {}
func (*MockLexerCallback) SongPositionPointer(beats uint16)                                {}
func (*MockLexerCallback) SongSelect(song uint8)                                           {}
func (*MockLexerCallback) Undefined1()                                                     {}
func (*MockLexerCallback) Undefined2()                                                     {}
func (*MockLexerCallback) TuneRequest()                                                    {}
func (*MockLexerCallback) TimingClock()                                                    {}
func (*MockLexerCallback) Undefined3()                                                     {}
func (*MockLexerCallback) Start()                                                          {}
func (*MockLexerCallback) Continue()                                                       {}
func (*MockLexerCallback) Stop()                                                           {}
func (*MockLexerCallback) Undefined4()                                                     {}
func (*MockLexerCallback) ActiveSensing()                                                  {}
func (*MockLexerCallback) Reset()                                                          {}
func (*MockLexerCallback) Done()                                                           {}

// A mock implementation of LexerCallback that counts each method call
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

	// Most recent values
	headerData  HeaderData
	chunkHeader ChunkHeader
}

func (cbk *CountingLexerCallback) Header(header HeaderData)                           { cbk.header++; cbk.headerData = header }
func (cbk *CountingLexerCallback) Track(header ChunkHeader)                           { cbk.track++; cbk.chunkHeader = header }
func (cbk *CountingLexerCallback) Began()                                             { cbk.began++ }
func (cbk *CountingLexerCallback) Finished()                                          { cbk.finished++ }
func (cbk *CountingLexerCallback) ErrorReading()                                      { cbk.errorReading++ }
func (cbk *CountingLexerCallback) ErrorOpeningFile()                                  { cbk.errorOpeningFile++ }
func (cbk *CountingLexerCallback) NoteOff(channel uint8, pitch uint8, velocity uint8) { cbk.noteOff++ }
func (cbk *CountingLexerCallback) NoteOn(channel uint8, pitch uint8, velocity uint8)  { cbk.noteOn++ }
func (cbk *CountingLexerCallback) PolyphonicAfterTouch(channel uint8, pitch uint8, pressure uint8) {
	cbk.polyphonicAfterTouch++
}
func (cbk *CountingLexerCallback) ControlChange(channel uint8, controller uint8, value uint8) {
	cbk.controlChange++
}
func (cbk *CountingLexerCallback) ProgramChange(channel uint8, program uint8) { cbk.programChange++ }
func (cbk *CountingLexerCallback) ChannelAfterTouch(channel uint8, value uint8) {
	cbk.channelAfterTouch++
}
func (cbk *CountingLexerCallback) PitchWheel(channel uint8, value uint16) { cbk.pitchWheel++ }
func (cbk *CountingLexerCallback) TimeCodeQuarter(messageType uint8, values uint8) {
	cbk.timeCodeQuarter++
}
func (cbk *CountingLexerCallback) SongPositionPointer(beats uint16) { cbk.songPositionPointer++ }
func (cbk *CountingLexerCallback) SongSelect(song uint8)            { cbk.songSelect++ }
func (cbk *CountingLexerCallback) Undefined1()                      { cbk.undefined1++ }
func (cbk *CountingLexerCallback) Undefined2()                      { cbk.undefined2++ }
func (cbk *CountingLexerCallback) TuneRequest()                     { cbk.tuneRequest++ }
func (cbk *CountingLexerCallback) TimingClock()                     { cbk.timingClock++ }
func (cbk *CountingLexerCallback) Undefined3()                      { cbk.undefined3++ }
func (cbk *CountingLexerCallback) Start()                           { cbk.start++ }
func (cbk *CountingLexerCallback) Continue()                        { cbk.continue_++ }
func (cbk *CountingLexerCallback) Stop()                            { cbk.stop++ }
func (cbk *CountingLexerCallback) Undefined4()                      { cbk.undefined4++ }
func (cbk *CountingLexerCallback) ActiveSensing()                   { cbk.activeSensing++ }
func (cbk *CountingLexerCallback) Reset()                           { cbk.reset++ }
func (cbk *CountingLexerCallback) Done()                            { cbk.done++ }

// MockReadSeeker is a mock Reader and Seeker. Constructed with data, behaves as a file reader.
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

// Seek sets the offset for the next Read or Write to offset, interpreted according to whence: 0 means relative to the origin of the file, 1 means relative to the current offset, and 2 means relative to the end. Seek returns the new offset and an Error, if any.
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
