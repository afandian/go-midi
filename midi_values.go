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

// Time code formats used in HeaderData.timeFormat
const (
	MetricalTimeFormat = iota
	TimeCodeTimeFormat = iota
)
