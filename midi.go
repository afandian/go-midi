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

// MidiLexer is a Standard Midi File Lexer.
// Pass this a Reader to a MIDI file and a callback that conforms to MidiLexerCallback 
// and it'll run over the file, calling events on the callback.
type MidiLexer struct {
	callback MidiLexerCallback
	input    io.Reader
}

// Construct a new MidiLexer
func NewMidiLexer(input io.Reader, callback MidiLexerCallback) *MidiLexer {
	return &MidiLexer{callback: callback, input: input}
}

// Lex starts the MidiLexer running.
func (lexer *MidiLexer) Lex() (error int) {
	if lexer.callback == nil {
		return NoCallback
	}

	if lexer.input == nil {
		return NoReader
	}

	return Ok
}

// next lexes the next item.
func (lexer *MidiLexer) next() (finished bool) {
	return true
}
