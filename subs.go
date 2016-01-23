/*

This file defines types to model subtitles and contains utilities.

*/

package srtgears

import (
	"time"
)

// Position of a subtitle.
type Pos int

// Possible position values of subtitles.
// Zero value of type Pos is PosNotSpecified.
const (
	PosNotSpecified Pos = iota
	BottomLeft
	Bottom
	BottomRight
	Left
	Center
	Right
	TopLeft
	Top
	TopRight
)

// Subtitle represents 1 subtitle, 1 displayable text (which may be broken into multiple lines).
type Subtitle struct {
	TimeIn  time.Duration // Timestamp when subtitle appears
	TimeOut time.Duration // Timestamp when subtitle disappears
	Lines   []string      // Lines of text to be displayed
	Pos     Pos           // Position where to display it
	Color   string        // Color of the text, HTML ARGB format
}

// DisplayDuration returns the duration for which the subtitle is visible.
func (s *Subtitle) DisplayDuration() time.Duration {
	return s.TimeOut - s.TimeIn
}

// HearingImpaired tells if the subtitle is hearing impaired.
func (s *Subtitle) HearingImpaired() bool {
	if len(s.Lines) == 0 {
		return false
	}
	lastLine := s.Lines[len(s.Lines)-1]
	first, last := s.Lines[0][0], lastLine[len(lastLine)-1]
	return first == '[' && last == ']' || first == '(' && last == ')'
}

// Shift shifts the subtitle with the specified delta.
func (s *Subtitle) Shift(delta time.Duration) {
	s.TimeIn += delta
	s.TimeOut += delta
}

// Scale scales the appearance timestamp of the subtitle.
// The duration for which subtitle is visible is not changed.
func (s *Subtitle) Scale(factor float64) {
	dispdur := s.DisplayDuration()
	s.TimeIn = time.Duration(float64(s.TimeIn) * factor)
	s.TimeOut = s.TimeIn + dispdur
}

// SubsPack represents subtitles of a movie,
// a collection of Subtitles and other meta info.
type SubsPack struct {
	Subs []*Subtitle
}

// Shift shifts all subtitles with the specified delta.
func (sp *SubsPack) Shift(delta time.Duration) {
	for _, s := range sp.Subs {
		s.Shift(delta)
	}
}

// Scale scales the timestamps of the subtitles.
// The duration for which subtitles are visible is not changed.
func (sp *SubsPack) Scale(factor float64) {
	for _, s := range sp.Subs {
		s.Scale(factor)
	}
}

// SetPos sets the position of all subtitles.
func (sp *SubsPack) SetPos(pos Pos) {
	for _, s := range sp.Subs {
		s.Pos = pos
	}
}

// RemoveHearingImpaired removes hearing impact-only subtitles
// (such as "[PHONE RINGING]" or "(phone ringing)").
func (sp *SubsPack) RemoveHearingImpaired() {
	for i := len(sp.Subs) - 1; i >= 0; i-- {
		if sp.Subs[i].HearingImpaired() {
			sp.Subs = append(sp.Subs[:i], sp.Subs[i+1:]...)
		}
	}
}
