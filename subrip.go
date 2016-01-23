/*

This file implements reading and writing the SubRip file format (*.srt).
It can parse *.srt files and create model from them.
And it also generate SubRip content from a model.

Format specifications:
https://en.wikipedia.org/wiki/SubRip
http://www.matroska.org/technical/specs/subtitles/srt.html

Unofficial extensions are also supported and used.

*/

package srtgears

import (
	"io"
)

// ParseSrt parses a SubRip stream (*.srt) and builds the model from it.
func ParseSrt(r io.Reader) (*SubsPack, error) {
	return nil, nil
}

//WriteSrt generates SubRip format.
func WriteSrt(w io.Writer, s *SubsPack) error {
	return nil
}
