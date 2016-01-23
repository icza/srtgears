/*

This file defines the Subtitle model type and its utility methods / transformations.

*/

package srtgears

import (
	"regexp"
	"time"
)

// Subtitle position type.
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

// RemoveHearingImpaired removes hearing impaired lines
// (such as "[PHONE RINGING]" or "(phone ringing)").
func (s *Subtitle) RemoveHearingImpaired() {
	// It may be just some (e.g. first) lines are hearing impaired.
	for i := len(s.Lines) - 1; i >= 0; i-- {
		line := s.Lines[i]
		first, last := line[0], line[len(line)-1]
		if first == '[' && last == ']' || first == '(' && last == ')' {
			s.Lines = append(s.Lines[:i], s.Lines[i+1:]...)
		}
	}
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

// Pattern used to remove HTML formatting
var htmlPattern = regexp.MustCompile(`<[^>]+>`)

// RemoveHTML removes HTML formatting.
func (s *Subtitle) RemoveHTML() {
	for i, v := range s.Lines {
		s.Lines[i] = htmlPattern.ReplaceAllString(v, "")
	}
}

// Pattern used to remove controls such as {\anX} (or {\aY}), {\pos(x,y)}.
var controlPattern = regexp.MustCompile(`^{\\[^}]*}`)

// RemoveControl removes controls such as {\anX} (or {\aY}), {\pos(x,y)}.
func (s *Subtitle) RemoveControl() {
	for i, v := range s.Lines {
		s.Lines[i] = controlPattern.ReplaceAllString(v, "")
	}
}
