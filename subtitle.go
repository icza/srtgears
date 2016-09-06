/*

This file defines the Subtitle model type and its utility methods / transformations.

*/

package srtgears

import (
	"regexp"
	"time"
)

// Pos is the subtitle position type.
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
	Color   string        // Color of the text, HTML RRGGBB format or a color name
}

// DisplayDuration returns the duration for which the subtitle is visible.
func (s *Subtitle) DisplayDuration() time.Duration {
	return s.TimeOut - s.TimeIn
}

// RemoveHI removes hearing impaired lines
// (such as "[PHONE RINGING]" or "(phone ringing)").
// Returns true if HI lines were present.
func (s *Subtitle) RemoveHI() (remove bool) {
	// It may be just some (e.g. first) lines are hearing impaired.
	for i := len(s.Lines) - 1; i >= 0; i-- {
		line := s.Lines[i]
		// Check without HTML formatting to recognize and remove these:
		// "<i>[PHONE RINGING]</i>"
		line = htmlPattern.ReplaceAllString(line, "")
		first, last := line[0], line[len(line)-1]
		if first == '[' && last == ']' || first == '(' && last == ')' {
			remove = true
			s.Lines = append(s.Lines[:i], s.Lines[i+1:]...)
		}
	}
	return
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

// Lengthen lengthens the display duration of the subtitle.
func (s *Subtitle) Lengthen(factor float64) {
	newDur := time.Duration(float64(s.DisplayDuration()) * factor)
	center := (s.TimeIn + s.TimeOut) / 2
	s.TimeIn = center - newDur/2
	if s.TimeIn < 0 { // Make sure it's not negative
		s.TimeIn = 0
	}
	s.TimeOut = s.TimeIn + newDur // specify relative to s.TimeIn as it may have been modified
}

// Pattern used to remove HTML formatting
var htmlPattern = regexp.MustCompile(`<[^>]+>`)

// RemoveHTML removes HTML formatting.
// Returns true if HTML formatting was present.
func (s *Subtitle) RemoveHTML() (removed bool) {
	for i, v := range s.Lines {
		s.Lines[i] = htmlPattern.ReplaceAllString(v, "")
		removed = removed || s.Lines[i] != v
	}
	// Color comes from HTML, so also zero it
	removed = removed || s.Color != ""
	s.Color = ""
	return
}

// Pattern used to remove controls such as {\anX} (or {\aY}), {\pos(x,y)}.
var controlPattern = regexp.MustCompile(`^{\\[^}]*}`)

// RemoveControl removes controls such as {\anX} (or {\aY}), {\pos(x,y)}.
// Returns true if controls were present.
func (s *Subtitle) RemoveControl() (removed bool) {
	for i, v := range s.Lines {
		s.Lines[i] = controlPattern.ReplaceAllString(v, "")
		removed = removed || s.Lines[i] != v
	}
	// Pos comes from control, so also zero it
	removed = removed || s.Pos != PosNotSpecified
	s.Pos = PosNotSpecified
	return
}
