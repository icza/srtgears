/*

This file implements reading and writing the SubRip file format (*.srt).
It can parse *.srt files and create model from them.
And it can also generate SubRip content from a model.

Format specifications:
https://en.wikipedia.org/wiki/SubRip
http://www.matroska.org/technical/specs/subtitles/srt.html

The parser is permissive, it tries to parse the input even if it does not conform to the specification.

Unofficial extensions are also supported and used.

An example SRT file:

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
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Mapping between *.srt pos to our model Pos for the {\anX} variant
var srtPosToModelPos = map[byte]Pos{
	'7': TopLeft, '8': Top, '9': TopRight,
	'4': Left, '5': Center, '6': Right,
	'1': BottomLeft, '2': Bottom, '3': BottomRight,
}

// Mapping between our model Pos to *.srt for the {\anX} variant
var modelPosToSrtPos = map[Pos]byte{
	TopLeft: '7', Top: '8', TopRight: '9',
	Left: '4', Center: '5', Right: '6',
	BottomLeft: '1', Bottom: '2', BottomRight: '3',
}

// ReadSrtFile reads and parses a SubRip file (*.srt) and builds the model from it.
func ReadSrtFile(name string) (sp *SubsPack, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	debugf("Reading from file: %s", name)
	return ReadSrtFrom(f)
}

// Regexp pattern to validate sequence number lines
var seqNumPattern = regexp.MustCompile(`^\s*\d+\s*$`)

// A starter font pattern and the color attribute value.
var starterFontPattern = regexp.MustCompile(`^\s*<\s*font\s+color\s*=\s*['"]?([^'"]*)['"]?\s*>`)

// </font> closing tag pattern.
var fontClosingPattern = regexp.MustCompile(`<\s*/\s*font\s*>`)

// ReadSrtFrom reads and parses a SubRip from an io.Reader (*.srt) and builds the model from it.
func ReadSrtFrom(r io.Reader) (sp *SubsPack, err error) {
	sp = &SubsPack{}
	scanner := bufio.NewScanner(r)
	phase := 0
	var s *Subtitle

	addSub := func() {
		// Post process
		if len(s.Lines) > 0 {
			// find position spec in first line (e.g. {\anX})
			if line := s.Lines[0]; strings.HasPrefix(line, `{\a`) {
				// 2 variants: {\anX} and {\aX}
				if len(line) >= 6 && line[3] == 'n' && line[5] == '}' {
					if p, ok := srtPosToModelPos[line[4]]; ok {
						s.Pos = p
						s.Lines[0] = line[6:] // Cut off pos spec from text
					}
				} else {
					// TODO other variant
				}
			}
			// Find if there is starter <font color="">
			if parts := starterFontPattern.FindStringSubmatch(s.Lines[0]); len(parts) > 0 {
				s.Color = parts[1]
				s.Lines[0] = s.Lines[0][len(parts[0]):] // cut off whole <font>
				// Find it's closing part
				for i, v := range s.Lines {
					if loc := fontClosingPattern.FindStringIndex(v); loc != nil {
						s.Lines[i] = v[:loc[0]] + v[loc[1]:]
						break
					}
				}
			}
		}
		sp.Subs, s = append(sp.Subs, s), nil
	}

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		if lineNum == 0 {
			// If BOM is present, strip it off. It's "\uFEFF", which is "\xef\xbb\xbf" in UTF-8
			if strings.HasPrefix(line, "\xef\xbb\xbf") {
				line = line[3:]
			}
		}
		lineNum++
		switch phase {
		case 0: // wanting sequence number, starting a new sub
			if line == "" {
				break // If multiple empty line separates, just ignore them
			}
			if Debug {
				if !seqNumPattern.MatchString(line) {
					debugf("Invalid sequence number in line %d: %s", lineNum, line)
				}
			}
			// discard seq#, we generate sequence numbers when writing
			s = &Subtitle{}
			phase++
		case 1: // wanting timestamps
			parseTimestamps(s, line, lineNum)
			phase++
		case 2: // wanting subtitle lines
			if line == "" {
				// End of subtitle, separator
				addSub()
				phase = 0
			} else {
				s.Lines = append(s.Lines, line)
			}
		}
	}
	if s != nil { // Append last subtitle if there is no empty line at the end of input
		addSub()
	}

	debugf("Loaded %d subtitles.", len(sp.Subs))

	sp.Sort()

	err = scanner.Err()
	return
}

// Regexp pattern to extract data from timestamp lines.
// Very permissive, for example also accepts this line:
//     dY 00:02:20.476--->   00:02:22,501X Y
var timestampsPattern = regexp.MustCompile(`(\d\d):(\d\d):(\d\d)[,\.](\d\d\d)\s*-+>\s*(\d\d):(\d\d):(\d\d)[,\.](\d\d\d)`)

//                                            0 0 :  0 0 :  0 0  ,     0 0 0    -->     0 0 :  0 0 :  0 0  ,     0 0 0

// parseTimestamps parses a timestamp line
func parseTimestamps(s *Subtitle, line string, lineNum int) {
	// Example: 00:02:20,476 --> 00:02:22,501
	parts := timestampsPattern.FindStringSubmatch(line)
	if len(parts) == 0 {
		// No match, invalid timestamp line
		debugf("Invalid timestamp in line %d: %s", lineNum, line)
		return
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
	s.TimeOut = time.Hour*get(5) + time.Minute*get(6) + time.Second*get(7) + time.Millisecond*get(8)

	if s.TimeOut <= s.TimeIn {
		debugf("Time1 >= Time2, text won't be visible in line %d: %s", lineNum, line)
	}
}

// WriteSrtFile generates SubRip format (*.srt) and writes it to a file.
func WriteSrtFile(name string, sp *SubsPack) (err error) {
	f, err := os.Create(name)
	if err != nil {
		return
	}
	defer f.Close()

	return WriteSrtTo(f, sp)
}

// WriteSrtTo generates SubRip format (*.srt) and writes it to an io.Writer.
func WriteSrtTo(w io.Writer, sp *SubsPack) error {
	wr := &writer{w: w}

	printTime := func(t time.Duration) {
		hour := t / time.Hour
		min := (t % time.Hour) / time.Minute
		sec := (t % time.Minute) / time.Second
		ms := (t % time.Second) / time.Millisecond
		wr.prf("%02d:%02d:%02d.%03d", hour, min, sec, ms)
	}

	for i, s := range sp.Subs {
		if wr.err != nil {
			break
		}

		// Sequence number
		wr.prn(i + 1)

		// Timestamps
		printTime(s.TimeIn)
		wr.pr(" --> ")
		printTime(s.TimeOut)
		wr.prn()

		// Texts
		for i, line := range s.Lines {
			if i == 0 && s.Pos != PosNotSpecified {
				wr.prf(`{\an%c}`, modelPosToSrtPos[s.Pos])
			}
			if s.Color != "" {
				// If there is color, wrap all lines into a <font>.
				if i == 0 { // This means opening in first line
					wr.prf(`<font color="%s">`, s.Color)
				}
				wr.pr(line)
				if i == len(s.Lines)-1 { // And closing in the last
					wr.pr("</font>")
				}
				wr.prn()
			} else {
				wr.prn(line)
			}
		}

		// Separator: empty line
		wr.prn()
	}

	return wr.err
}
