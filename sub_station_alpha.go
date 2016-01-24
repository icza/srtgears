/*

This file implements writing the Sub Station Alpha file format (*.ssa).
It can generate SubRip content from a model.

Format specifications:
https://en.wikipedia.org/wiki/SubStation_Alpha
http://www.matroska.org/technical/specs/subtitles/ssa.html

An example SSA file:

	[Script Info]
	; This is a Sub Station Alpha v4 script.
	; For Sub Station Alpha info and downloads,
	; go to http://www.eswat.demon.co.uk/
	Title: Neon Genesis Evangelion - Episode 26 (neutral Spanish)
	Original Script: RoRo
	Script Updated By: version 2.8.01
	ScriptType: v4.00
	Collisions: Normal
	PlayResY: 600
	PlayDepth: 0
	Timer: 100,0000

	[V4 Styles]
	Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, TertiaryColour, BackColour, Bold, Italic, BorderStyle, Outline, Shadow,
	   Alignment, MarginL, MarginR, MarginV, AlphaLevel, Encoding
	Style: DefaultVCD, Arial,28,11861244,11861244,11861244,-2147483640,-1,0,1,1,2,2,30,30,30,0,0

	[Events]
	Format: Marked, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
	Dialogue: Marked=0,0:00:01.18,0:00:06.85,DefaultVCD, NTP,0000,0000,0000,,{\pos(400,570)}Like an angel with pity on nobody

*/

package srtgears

import (
	"io"
	"os"
)

// Mapping between our model Pos to *.ssa Alignment
var modelPosToSsaPos = map[Pos]int{
	TopLeft: 5, Top: 6, TopRight: 7,
	Left: 9, Center: 10, Right: 11,
	BottomLeft: 1, Bottom: 2, BottomRight: 3,
}

// WriteSsaFile generates Sub Station Alpha format and writes it to a file.
func WriteSsaFile(name string, sp *SubsPack) (err error) {
	f, err := os.Create(name)
	if err != nil {
		return
	}
	defer f.Close()

	return WriteSsaTo(f, sp)
}

// WriteSsaTo generates Sub Station Alpha format and writes it to an io.Writer.
func WriteSsaTo(w io.Writer, sp *SubsPack) (err error) {
	// TODO
	return
}
