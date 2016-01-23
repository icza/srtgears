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
