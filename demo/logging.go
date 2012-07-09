package main

import (
	"fmt"
	"midi"
	"os"
)

type LoggingLexerCallback struct{}

func (cbk LoggingLexerCallback) Header(header midi.HeaderData) { fmt.Println("Header", header) }
func (cbk LoggingLexerCallback) Track(header midi.ChunkHeader) { fmt.Println("Track", header) }
func (cbk LoggingLexerCallback) Began()                        { fmt.Println("Began") }
func (cbk LoggingLexerCallback) Finished()                     { fmt.Println("Finished") }
func (cbk LoggingLexerCallback) ErrorReading()                 { fmt.Println("ErrorReading") }
func (cbk LoggingLexerCallback) ErrorOpeningFile()             { fmt.Println("ErrorOpeningFile") }
func (cbk LoggingLexerCallback) Tempo(bpm uint32, microsecondsPerCrotchet uint32, time uint32) {
	fmt.Println("tempo", bpm, "bpm")
}
func (cbk LoggingLexerCallback) NoteOff(channel uint8, pitch uint8, velocity uint8, time uint32) {
	// fmt.Println("NoteOff", channel, pitch, velocity, time)

	for i := uint32(0); i < time/100; i++ {
		fmt.Println("")
	}

	for i := uint8(0); i < pitch; i++ {
		fmt.Print(" ")
	}
	fmt.Println("x")

}
func (cbk LoggingLexerCallback) NoteOn(channel uint8, pitch uint8, velocity uint8, time uint32) {
	// fmt.Println("NoteOn", channel, pitch, velocity, time)

	for i := uint32(0); i < time/100; i++ {
		fmt.Println("")
	}

	for i := uint8(0); i < pitch; i++ {
		fmt.Print(" ")
	}
	fmt.Println("*")
}
func (cbk LoggingLexerCallback) PolyphonicAfterTouch(channel uint8, pitch uint8, pressure uint8, time uint32) {
	fmt.Println("PolyphonicAfterTouch", channel, pitch, pressure, time)
}
func (cbk LoggingLexerCallback) ControlChange(channel uint8, controller uint8, value uint8, time uint32) {
	fmt.Println("ControlChange", channel, controller, value, time)
}
func (cbk LoggingLexerCallback) ProgramChange(channel uint8, program uint8, time uint32) {
	fmt.Println("ProgramChange", channel, program, time)
}
func (cbk LoggingLexerCallback) ChannelAfterTouch(channel uint8, value uint8, time uint32) {
	fmt.Println("ChannelAfterTouch", channel, value, time)
}
func (cbk LoggingLexerCallback) PitchWheel(channel uint8, value int16, absValue uint16, time uint32) {
	fmt.Println("PitchWheel", channel, value, absValue, time)
}
func (cbk LoggingLexerCallback) TimeCodeQuarter(messageType uint8, values uint8, time uint32) {
	fmt.Println("TimeCodeQuarter", messageType, values, time)
}
func (cbk LoggingLexerCallback) SongPositionPointer(beats uint16, time uint32) {
	fmt.Println("SongPositionPointer", beats, time)
}
func (cbk LoggingLexerCallback) SongSelect(song uint8, time uint32) {
	fmt.Println("SongSelect", song, time)
}
func (cbk LoggingLexerCallback) Undefined1(time uint32)    { fmt.Println("Undefined1", time) }
func (cbk LoggingLexerCallback) Undefined2(time uint32)    { fmt.Println("Undefined2", time) }
func (cbk LoggingLexerCallback) TuneRequest(time uint32)   { fmt.Println("TuneRequest", time) }
func (cbk LoggingLexerCallback) TimingClock(time uint32)   { fmt.Println("TimingClock", time) }
func (cbk LoggingLexerCallback) Undefined3(time uint32)    { fmt.Println("Undefined3", time) }
func (cbk LoggingLexerCallback) Start(time uint32)         { fmt.Println("Start", time) }
func (cbk LoggingLexerCallback) Continue(time uint32)      { fmt.Println("Continue", time) }
func (cbk LoggingLexerCallback) Stop(time uint32)          { fmt.Println("Stop", time) }
func (cbk LoggingLexerCallback) Undefined4(time uint32)    { fmt.Println("Undefined4", time) }
func (cbk LoggingLexerCallback) ActiveSensing(time uint32) { fmt.Println("ActiveSensing", time) }
func (cbk LoggingLexerCallback) Reset(time uint32)         { fmt.Println("Reset", time) }
func (cbk LoggingLexerCallback) Done(time uint32)          { fmt.Println("Done", time) }
func (cbk LoggingLexerCallback) SequenceNumber(channel uint8, number uint16, numberGiven bool, time uint32) {
	fmt.Println("SequenceNumber", channel, number, numberGiven, time)
}
func (cbk LoggingLexerCallback) Text(channel uint8, text string, time uint32) {
	fmt.Println("Text", channel, text, time)
}
func (cbk LoggingLexerCallback) CopyrightText(channel uint8, text string, time uint32) {
	fmt.Println("CopyrightText", channel, text, time)
}
func (cbk LoggingLexerCallback) SequenceName(channel uint8, text string, time uint32) {
	fmt.Println("SequenceName", channel, text, time)
}
func (cbk LoggingLexerCallback) TrackInstrumentName(channel uint8, text string, time uint32) {
	fmt.Println("TrackInstrumentName", channel, text, time)
}
func (cbk LoggingLexerCallback) LyricText(channel uint8, text string, time uint32) {
	fmt.Println("LyricText", channel, text, time)
}
func (cbk LoggingLexerCallback) MarkerText(channel uint8, text string, time uint32) {
	fmt.Println("MarkerText", channel, text, time)
}
func (cbk LoggingLexerCallback) CuePointText(channel uint8, text string, time uint32) {
	fmt.Println("CuePointText", channel, text, time)
}
func (cbk LoggingLexerCallback) EndOfTrack(channel uint8, time uint32) {
	fmt.Println("EndOfTrack", channel, time)
}
func (cbk LoggingLexerCallback) TimeSignature(numerator uint8, denomenator uint8, clocksPerClick uint8, demiSemiQuaverPerQuarter uint8, time uint32) {
	fmt.Println("TimeSignature", numerator, denomenator, clocksPerClick, demiSemiQuaverPerQuarter, time)
}

func main() {
	fmt.Println("Logging Midi")
	var callback LoggingLexerCallback
	// var file, err = os.Open("/Users/joe/Downloads/102891.mid")
	var file, err = os.Open("/Users/joe/Downloads/HOTELCAL.MID")
	if err != nil {
		fmt.Println(err)
		return
	}
	lexer := midi.NewMidiLexer(file, callback)
	lexer.Lex()
}
