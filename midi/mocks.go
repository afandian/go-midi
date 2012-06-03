// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

// Mocks implementations for testing.

package midi

// A mock implementation of MockLexerCallback that does nothing
type MockLexerCallback struct {}
func (* MockLexerCallback) Began() {}
func (* MockLexerCallback) Finished() {}
func (* MockLexerCallback) ErrorReading() {}
func (* MockLexerCallback) ErrorOpeningFile() {}
func (* MockLexerCallback) NoteOff(channel uint8, pitch uint8, velocity uint8) {}
func (* MockLexerCallback) NoteOn(channel uint8, pitch uint8, velocity uint8) {}
func (* MockLexerCallback) PolyphonicAfterTouch(channel uint8, pitch uint8, pressure uint8) {}
func (* MockLexerCallback) ControlChange(channel uint8, controller uint8, value uint8) {}
func (* MockLexerCallback) ProgramChange(channel uint8, program uint8) {}
func (* MockLexerCallback) ChannelAfterTouch(channel uint8, value uint8) {}
func (* MockLexerCallback) PitchWheel(channel uint8, value uint16) {}
func (* MockLexerCallback) TimeCodeQuarter(messageType uint8, values uint8) {}
func (* MockLexerCallback) SongPositionPointer(beats uint16) {}
func (* MockLexerCallback) SongSelect(song uint8) {}
func (* MockLexerCallback) Undefined1() {}
func (* MockLexerCallback) Undefined2() {}
func (* MockLexerCallback) TuneRequest() {}
func (* MockLexerCallback) TimingClock() {}
func (* MockLexerCallback) Undefined3() {}
func (* MockLexerCallback) Start() {}
func (* MockLexerCallback) Continue() {}
func (* MockLexerCallback) Stop() {}
func (* MockLexerCallback) Undefined4() {}
func (* MockLexerCallback) ActiveSensing() {}
func (* MockLexerCallback) Reset() {}
func (* MockLexerCallback) Done() {}

// MockReader is a mock Reader. Constructed with data, behaves as a file reader.
type MockReader struct {
	data *[]byte
	position int
}

// NewMockReader creates a new MockReader object backed by the given byte array data.
func NewMockReader(data *[]byte) *MockReader {
	return &MockReader{data: data}
}

// Read fills the given buffer, returning the number of bytes and an error.
func (reader *MockReader) Read(p []byte) (n int, err error) {
	var amount = len(p)
	var maxAmount = len(*reader.data) - reader.position

	// Don't read past the end
	if amount > maxAmount {
		amount = maxAmount
	}

	copy(p, (*reader.data)[reader.position:reader.position+amount])
	reader.position += amount
	return amount, nil
}