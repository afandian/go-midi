package main

import (
	"midi"
	"os"
)

type LoggingLexerCallback struct{}

func (cbk LoggingLexerCallback) Header(header midi.HeaderData) {}
func (cbk LoggingLexerCallback) Track(header midi.ChunkHeader) {}
func (cbk LoggingLexerCallback) Began()                        {}
func (cbk LoggingLexerCallback) Finished()                     {}
func (cbk LoggingLexerCallback) ErrorReading()                 {}
func (cbk LoggingLexerCallback) ErrorOpeningFile()             {}
func (cbk LoggingLexerCallback) Tempo(bpm uint32, microsecondsPerCrotchet uint32, time uint32) {

}
func (cbk LoggingLexerCallback) NoteOff(channel uint8, pitch uint8, velocity uint8, time uint32) {

}
func (cbk LoggingLexerCallback) NoteOn(channel uint8, pitch uint8, velocity uint8, time uint32) {
}
func (cbk LoggingLexerCallback) PolyphonicAfterTouch(channel uint8, pitch uint8, pressure uint8, time uint32) {

}
func (cbk LoggingLexerCallback) ControlChange(channel uint8, controller uint8, value uint8, time uint32) {

}
func (cbk LoggingLexerCallback) ProgramChange(channel uint8, program uint8, time uint32) {

}
func (cbk LoggingLexerCallback) ChannelAfterTouch(channel uint8, value uint8, time uint32) {

}
func (cbk LoggingLexerCallback) PitchWheel(channel uint8, value int16, absValue uint16, time uint32) {

}
func (cbk LoggingLexerCallback) TimeCodeQuarter(messageType uint8, values uint8, time uint32) {

}
func (cbk LoggingLexerCallback) SongPositionPointer(beats uint16, time uint32) {

}
func (cbk LoggingLexerCallback) SongSelect(song uint8, time uint32) {

}
func (cbk LoggingLexerCallback) Undefined1(time uint32)    {}
func (cbk LoggingLexerCallback) Undefined2(time uint32)    {}
func (cbk LoggingLexerCallback) TuneRequest(time uint32)   {}
func (cbk LoggingLexerCallback) TimingClock(time uint32)   {}
func (cbk LoggingLexerCallback) Undefined3(time uint32)    {}
func (cbk LoggingLexerCallback) Start(time uint32)         {}
func (cbk LoggingLexerCallback) Continue(time uint32)      {}
func (cbk LoggingLexerCallback) Stop(time uint32)          {}
func (cbk LoggingLexerCallback) Undefined4(time uint32)    {}
func (cbk LoggingLexerCallback) ActiveSensing(time uint32) {}
func (cbk LoggingLexerCallback) Reset(time uint32)         {}
func (cbk LoggingLexerCallback) Done(time uint32)          {}
func (cbk LoggingLexerCallback) SequenceNumber(channel uint8, number uint16, numberGiven bool, time uint32) {

}
func (cbk LoggingLexerCallback) Text(channel uint8, text string, time uint32) {

}
func (cbk LoggingLexerCallback) CopyrightText(channel uint8, text string, time uint32) {

}
func (cbk LoggingLexerCallback) SequenceName(channel uint8, text string, time uint32) {

}
func (cbk LoggingLexerCallback) TrackInstrumentName(channel uint8, text string, time uint32) {

}
func (cbk LoggingLexerCallback) LyricText(channel uint8, text string, time uint32) {

}
func (cbk LoggingLexerCallback) MarkerText(channel uint8, text string, time uint32) {

}
func (cbk LoggingLexerCallback) CuePointText(channel uint8, text string, time uint32) {

}
func (cbk LoggingLexerCallback) EndOfTrack(channel uint8, time uint32) {

}
func (cbk LoggingLexerCallback) TimeSignature(numerator uint8, denomenator uint8, clocksPerClick uint8, demiSemiQuaverPerQuarter uint8, time uint32) {

}
func (cbk LoggingLexerCallback) KeySignature(key midi.ScaleDegree, mode midi.KeySignatureMode, sharpsOrFlats int8) {

}

func main() {

	var callback LoggingLexerCallback
	// loc := "/Users/joe/Downloads/102891.mid"
	// loc := "/Users/joe/Downloads/HOTELCAL.MID"

	loc := "/Users/joe/personal/backup/home/Websites/close site5/ttf/public_html_old-28-apr-2010/temporaryImages/1018.mid"
	var file, err = os.Open(loc)
	if err != nil {
		return
	}
	lexer := midi.NewMidiLexer(file, callback)
	lexer.Lex()
}
