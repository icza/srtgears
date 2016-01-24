/*

This file defines the SubsPack model type and its utility methods / transformations.

*/

package srtgears

import (
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

// SubsPack represents subtitles of a movie,
// a collection of Subtitles and other meta info.
type SubsPack struct {
	Subs []*Subtitle
}

// Type that implements sorting
type SortSubtitles []*Subtitle

func (s SortSubtitles) Len() int           { return len(s) }
func (s SortSubtitles) Less(i, j int) bool { return s[i].TimeIn < s[j].TimeIn }
func (s SortSubtitles) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Sort sorts the subtitles by appearance timestamp.
func (sp *SubsPack) Sort() {
	sort.Sort(SortSubtitles(sp.Subs))
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

// SetColor sets the color of all subtitles.
func (sp *SubsPack) SetColor(color string) {
	for _, s := range sp.Subs {
		s.Color = color
	}
}

// RemoveHI removes hearing impaired lines from subtitles
// (such as "[PHONE RINGING]" or "(phone ringing)").
func (sp *SubsPack) RemoveHI() {
	for i := len(sp.Subs) - 1; i >= 0; i-- {
		s := sp.Subs[i]
		s.RemoveHI()
		if len(s.Lines) == 0 {
			// Can be removed completely
			sp.Subs = append(sp.Subs[:i], sp.Subs[i+1:]...)
		}
	}
}

// Concatenate concatenates another SubsPack to this.
// Subtitles are not copied, only their addresses are appended to ours.
//
// In order to get correct timing for the concatenated 2nd part,
// timestamps of the concatenated subtitles have to be shifted
// with the start time of the 2nd part of the movie
// (which is usually the length of the first part).
//
// Useful if movie is present in 1 part but there are 2 subtitles for 2 parts.
func (sp *SubsPack) Concatenate(sp2 *SubsPack, secPartStart time.Duration) {
	sp2.Shift(secPartStart)
	sp.Subs = append(sp.Subs, sp2.Subs...)

	// There might be overlapping between the 2 parts
	// (e.g. 2nd part repeats the last minute of the first part),
	// so just to be well-behaved:
	sp.Sort()
}

// Merge merges another SubsPack into this to create a "dual subtitle".
// Subtitles are not copied, only their addresses are merged to ours.
//
// Useful if 2 different subtitles are to be displayed at the same time, e.g. 2 different languages.
func (sp *SubsPack) Merge(sp2 *SubsPack) {
	// Put our subtitles to bottom:
	sp.SetPos(PosNotSpecified)

	// Put other subtitles to the top:
	sp.SetPos(Top)

	// Append:
	sp.Subs = append(sp.Subs, sp2.Subs...)

	// And sort
	sp.Sort()
}

// Split splits this SubsPack into 2 at the specified time.
// Subtitles before the split time will remain in this, subtitles after the split time
// will be added to a new SubsPack that is returned.
//
// Useful to create 2 subtitles if movie is present in 2 parts but subtitle is for one.
func (sp *SubsPack) Split(at time.Duration) (sp2 *SubsPack) {
	sp2 = &SubsPack{}

	subs := sp.Subs
	idx := sort.Search(len(subs), func(i int) bool {
		return subs[i].TimeIn >= at
	})

	sp2.Subs = make([]*Subtitle, len(subs)-idx)
	copy(sp2.Subs, subs[idx:])
	sp.Subs = subs[:idx]

	// Shift splitted subs:
	sp2.Shift(-at)

	return
}

// RemoveHTML removes HTML formatting from all subtitles.
func (sp *SubsPack) RemoveHTML() {
	for _, s := range sp.Subs {
		s.RemoveHTML()
	}
}

// RemoveControl removes controls such as {\anX} (or {\aY}), {\pos(x,y)} from all subtitles.
func (sp *SubsPack) RemoveControl() {
	for _, s := range sp.Subs {
		s.RemoveControl()
	}
}

// Lengthen lenthens the display duration of all subtitles.
func (sp *SubsPack) Lengthen(factor float64) {
	for _, s := range sp.Subs {
		s.Lengthen(factor)
	}
}

// Statistics that can be gathered from a SubsPack.
type SubsStats struct {
	Subs                      int           // Total number of subs.
	Lines                     int           // # of lines
	AvgLinesPerSub            float64       // Avg lines per sub
	Chars                     int           // # of characters (spaces included)
	CharsNoSpace              int           // # of characters (without spaces)
	AvgCharsPerLine           float64       // Avg chars (no space) per line
	Words                     int           // # of words
	AvgWordsPerLine           float64       // Avg words per line
	AvgCharsPerWord           float64       // Avg chars per word
	TotalDispDur              time.Duration // Total subtitle display time
	SubVisibRatio             float64       // Subtitle visible ratio (compared to total time)
	AvgDispDurPerNonSpaceChar time.Duration // Avg. display duration per non-space char
	HTMLs                     int           // # of subs having HTML formatting
	Controls                  int           // # of subs having controls
	HIs                       int           // # of subs having hearing impaired lines
}

// Stats analyzes the subtitle pack and returns various statistics.
// Subtitles will be modified.
func (sp *SubsPack) Stats() *SubsStats {
	ss := SubsStats{
		Subs: len(sp.Subs),
	}

	for _, s := range sp.Subs {
		ss.TotalDispDur += s.DisplayDuration()
		ss.Lines += len(s.Lines)

		if s.RemoveControl() {
			ss.Controls++
		}
		if s.RemoveHTML() {
			ss.HTMLs++
		}

		for _, v := range s.Lines {
			ss.Chars += utf8.RuneCountInString(v)
			fields := strings.Fields(v)
			ss.Words += len(fields)
			for _, v2 := range fields {
				ss.CharsNoSpace += utf8.RuneCountInString(v2)
			}
		}

		if s.RemoveHI() {
			ss.HIs++
		}
	}

	if len(sp.Subs) > 0 {
		if last := sp.Subs[len(sp.Subs)-1].TimeOut; last != 0 {
			ss.SubVisibRatio = float64(ss.TotalDispDur) / float64(last)
		}
	}

	ss.AvgLinesPerSub = float64(ss.Lines) / float64(ss.Subs)
	ss.AvgCharsPerLine = float64(ss.CharsNoSpace) / float64(ss.Lines)
	ss.AvgWordsPerLine = float64(ss.Words) / float64(ss.Lines)
	ss.AvgCharsPerWord = float64(ss.CharsNoSpace) / float64(ss.Words)
	if ss.CharsNoSpace > 0 {
		ss.AvgDispDurPerNonSpaceChar = ss.TotalDispDur / time.Duration(ss.CharsNoSpace)
	}
	return &ss
}
