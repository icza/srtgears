/*

This file defines types to model subtitles and contains utilities.

*/

package srtgears

import (
	"time"
)

// Position of a subtitle.
type Pos int

// Possible values of subtitles
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

// Subtitle represents 1 subtitle, 1 displayable text.
type Subtitle struct {
	TimeIn  time.Duration // Timestamp when subtitle appears
	TimeOut time.Duration // Timestamp when subtitle disappears
	Lines   []string      // Lines of text to be displayed
	Pos     Pos           // Position where to display it
	Color   string        // Color of the text, HTML ARGB format
}

// SubsPack represents subtitles of a movie,
// a collection of Subtitles and other meta info.
type SubsPack struct {
	Subs []*Subtitle
}
