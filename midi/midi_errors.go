// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

// Data structures

package midi

// Error codes for Lexer
const (
	Ok         = 0x01
	NoCallback = 0x01 << 1
	NoReader   = 0x01 << 2
)

// A load of Errors and single values for convenience.

type VarLengthNotFoundError struct{}

func (e VarLengthNotFoundError) Error() string {
	return "Variable length value not found where expected."
}

type UnexpectedEndOfFileError struct{}

func (e UnexpectedEndOfFileError) Error() string {
	return "Unexpected End of File found."
}

var UnexpectedEndOfFile = UnexpectedEndOfFileError{}

type UnsupportedSmfFormatError struct{}

func (e UnsupportedSmfFormatError) Error() string {
	return "The SMF format was not expected."
}

var UnsupportedSmfFormat = UnsupportedSmfFormatError{}
