/*

This file implements reading and writing the SubRip file format (*.srt).
It can parse *.srt files and create model from them.
And it also generate SubRip content from a model.

Format specifications:
https://en.wikipedia.org/wiki/SubRip
http://www.matroska.org/technical/specs/subtitles/srt.html

The parser is permissive, it tries to parse the input even if it does not conform to the specification.

Unofficial extensions are also supported and used.

Example of an SRT file:

    1
    00:02:17,440 --> 00:02:20,375
    Senator, we're making
    our final approach into Coruscant.

    2
    00:02:20,476 --> 00:02:22,501
    Very good, Lieutenant.

*/

package srtgears

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"time"
)

// ReadSrt reads and parses a SubRip stream (*.srt) and builds the model from it.
func ReadSrt(r io.Reader) (sp *SubsPack, err error) {
	sp = &SubsPack{}

	scanner := bufio.NewScanner(r)

	phase := 0

	var s *Subtitle
	for scanner.Scan() {
		line := scanner.Text()

		switch phase {
		case 0: // wanting sequence number, starting a new sub
			// discard seq#, we generate sequence numbres when writing
			s = &Subtitle{}
			phase++
		case 1: // wanting timestamps
			parseTimestamps(s, line)
			phase++
		case 2: // wanting subtitle lines
			if line == "" {
				// End of subtitle, separator
				sp.Subs, s = append(sp.Subs, s), nil
				phase = 0
			} else {
				s.Lines = append(s.Lines, line)
			}
		}
	}

	err = scanner.Err()
	return
}

// Regexp pattern to extract data from timestamp lines.
// Very permissive, for example also accepts this line:
//     dY 00:02:20.476--->   00:02:22,501X Y
var timestampsPattern = regexp.MustCompile(`(\d\d):(\d\d):(\d\d)[,\.](\d\d\d)\s*-+>\s*(\d\d):(\d\d):(\d\d)[,\.](\d\d\d)`)

//                                            0 0 :  0 0 :  0 0  ,     0 0 0    -->     0 0 :  0 0 :  0 0  ,     0 0 0

// parseTimestamps parses a timestamp line
func parseTimestamps(s *Subtitle, line string) {
	// Example: 00:02:20,476 --> 00:02:22,501
	parts := timestampsPattern.FindStringSubmatch(line)
	if len(parts) == 0 {
		return // No match, invalid timestamp line
	}

	get := func(idx int) time.Duration {
		n, err := strconv.ParseInt(parts[idx], 10, 64)
		if err != nil {
			panic(err) // This shouldn't happen as only digits are matched.
		}
		return time.Duration(n)
	}

	// First part is the complete match
	s.TimeIn = time.Hour*get(1) + time.Minute*get(2) + time.Second*get(3) + time.Millisecond*get(4)
	s.TimeIn = time.Hour*get(5) + time.Minute*get(6) + time.Second*get(7) + time.Millisecond*get(8)
}

//WriteSrt generates SubRip format.
func WriteSrt(w io.Writer, s *SubsPack) error {
	return nil
}
