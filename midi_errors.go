// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

/*
 * Errors that may arise during parsing.
 * The LexerCallback may recieve any of these.
 */

package midi

// A load of Errors and single values for convenience.

type UnexpectedEventLengthError struct {
	message string
}

func (e UnexpectedEventLengthError) Error() string {
	return e.message
}

type NoCallbackError struct{}

func (e NoCallbackError) Error() string {
	return "No callback supplied"
}

var NoCallback = NoCallbackError{}

type NoReadSeekerError struct{}

func (e NoReadSeekerError) Error() string {
	return "No ReadSeeker input supplied"
}

var NoReadSeeker = NoReadSeekerError{}

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

type ExpectedMthdError struct{}

func (e ExpectedMthdError) Error() string {
	return "Expected SMF Midi header."
}

var ExpectedMthd = ExpectedMthdError{}

type BadSizeChunkError struct{}

func (e BadSizeChunkError) Error() string {
	return "Chunk was an unexpected size."
}

var BadSizeChunk = BadSizeChunkError{}
