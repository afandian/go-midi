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
	"testing"
)


// MidiLexer should throw error for null callback or input
func TestLexerShouldComplainNullArgs(t *testing.T) {
	var lexer *MidiLexer

	var mockLexerCallback = new(MockLexerCallback)
	var mockReadSeeker = NewMockReadSeeker(&[]byte{})
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

