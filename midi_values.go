// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

// Constants and values.

package midi

// SMF format
const (
	SingleMultiTrackChannel = 0
	SimultaneousTracks      = 1
	SequentialTracks        = 2
)

// Time code formats used in HeaderData.TimeFormat
const (
	MetricalTimeFormat = iota
	TimeCodeTimeFormat = iota
)

// Supplied to KeySignature
const (
	MajorMode = 0
	MinorMode = 1
)

type KeySignatureMode uint8

const (
	DegreeC  = 0
	DegreeCs = 1
	DegreeDf = DegreeCs
	DegreeD  = 2
	DegreeDs = 3
	DegreeEf = DegreeDs
	DegreeE  = 4
	DegreeF  = 5
	DegreeFs = 6
	DegreeGf = DegreeFs
	DegreeG  = 7
	DegreeGs = 8
	DegreeAf = DegreeGs
	DegreeA  = 9
	DegreeAs = 10
	DegreeBf = DegreeAs
	DegreeB  = 11
	DegreeCf = DegreeB
)

type ScaleDegree uint8
