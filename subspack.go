/*

This file defines the SubsPack model type and its utility methods / transformations.

*/

package srtgears

import (
	"sort"
	"time"
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

// RemoveHearingImpaired removes hearing impaired subtitles
// (such as "[PHONE RINGING]" or "(phone ringing)").
func (sp *SubsPack) RemoveHearingImpaired() {
	for i := len(sp.Subs) - 1; i >= 0; i-- {
		if sp.Subs[i].HearingImpaired() {
			sp.Subs = append(sp.Subs[:i], sp.Subs[i+1:]...)
		}
	}
}

// Concatenate concatenates another SubsPack into this.
// Subtitles are not copied, only their address is appended to ours.
//
// In order to get correct timing for the concatenated 2nd part,
// timestamps of the concatenated subtitles have to be shifted
// with the start time of the 2nd part of the movie
// (which is usually the length of the first part).
func (sp *SubsPack) Concatenate(sp2 *SubsPack, secPartStart time.Duration) {
	sp2.Shift(secPartStart)
	sp.Subs = append(sp.Subs, sp2.Subs...)

	// there might be overlapping between the 2 parts
	// (e.g. 2nd part repeats the last minute of the end of the first part),
	// so just to be well-behaved:
	sp.Sort()
}
