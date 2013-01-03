// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

/*
 * Tests for lexer.
 * Check that the state transitions work fine and that the lexer can cope with real streams of MIDI data.
 */

package midi

import (
	"testing"
)

func keyTest(sharpsFlats int, mode KeySignatureMode, expectedDegree ScaleDegree, expectedMode KeySignatureMode, message string, t *testing.T) {
	gotKey, gotMode := keySignatureFromSharpsOrFlats(sharpsFlats, uint(mode))
	if !((gotKey == expectedDegree) && (gotMode == expectedMode)) {
		t.Fatal("Fail", message, "key", gotKey, "e", expectedDegree, "mode", gotMode, "e", expectedMode)

	}
}

// Test lots of key signatures.
// Might as well test every possible input.
func TestKeySignatureTest(t *testing.T) {
	keyTest(
		0, MajorMode,
		DegreeC, MajorMode,
		"C Major", t,
	)

	keyTest(
		0, MinorMode,
		DegreeA, MinorMode,
		"A Minor", t,
	)

	// Sharps
	keyTest(
		1, MajorMode,
		DegreeG, MajorMode,
		"G Major", t,
	)

	keyTest(
		1, MinorMode,
		DegreeE, MinorMode,
		"E Minor", t,
	)

	keyTest(
		2, MajorMode,
		DegreeD, MajorMode,
		"D Major", t,
	)

	keyTest(
		2, MinorMode,
		DegreeB, MinorMode,
		"B Minor", t,
	)

	keyTest(
		3, MajorMode,
		DegreeA, MajorMode,
		"A Major", t,
	)

	keyTest(
		3, MinorMode,
		DegreeFs, MinorMode,
		"Fs Minor", t,
	)

	keyTest(
		4, MajorMode,
		DegreeE, MajorMode,
		"E Major", t,
	)

	keyTest(
		4, MinorMode,
		DegreeCs, MinorMode,
		"Cs Minor", t,
	)

	keyTest(
		5, MajorMode,
		DegreeB, MajorMode,
		"B Major", t,
	)

	keyTest(
		5, MinorMode,
		DegreeGs, MinorMode,
		"Gs Minor", t,
	)

	keyTest(
		6, MajorMode,
		DegreeFs, MajorMode,
		"Fs Major", t,
	)

	keyTest(
		6, MinorMode,
		DegreeDs, MinorMode,
		"Ds Minor", t,
	)

	keyTest(
		7, MajorMode,
		DegreeCs, MajorMode,
		"Cs Major", t,
	)

	keyTest(
		7, MinorMode,
		DegreeAs, MinorMode,
		"As Minor", t,
	)

	// Flats
	keyTest(
		-1, MajorMode,
		DegreeF, MajorMode,
		"F Major", t,
	)

	keyTest(
		-1, MinorMode,
		DegreeD, MinorMode,
		"D Minor", t,
	)

	keyTest(
		-2, MajorMode,
		DegreeBf, MajorMode,
		"Bf Major", t,
	)

	keyTest(
		-2, MinorMode,
		DegreeG, MinorMode,
		"G Minor", t,
	)

	keyTest(
		-3, MajorMode,
		DegreeEf, MajorMode,
		"Ef Major", t,
	)

	keyTest(
		-3, MinorMode,
		DegreeC, MinorMode,
		"C Minor", t,
	)

	keyTest(
		-4, MajorMode,
		DegreeAf, MajorMode,
		"Af Major", t,
	)

	keyTest(
		-4, MinorMode,
		DegreeF, MinorMode,
		"F Minor", t,
	)

	keyTest(
		-5, MajorMode,
		DegreeDf, MajorMode,
		"Df Major", t,
	)

	keyTest(
		-5, MinorMode,
		DegreeBf, MinorMode,
		"Bf Minor", t,
	)

	keyTest(
		-6, MajorMode,
		DegreeGf, MajorMode,
		"Gf Major", t,
	)

	keyTest(
		-6, MinorMode,
		DegreeEf, MinorMode,
		"Ef Minor", t,
	)

	keyTest(
		-7, MajorMode,
		DegreeCf, MajorMode,
		"Cf Major", t,
	)

	keyTest(
		-7, MinorMode,
		DegreeAf, MinorMode,
		"Af Minor", t,
	)
}
