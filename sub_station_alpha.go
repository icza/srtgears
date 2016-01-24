/*

This file implements writing the Sub Station Alpha file format (*.ssa).
It can generate SubRip content from a model.

Format specifications:
https://en.wikipedia.org/wiki/SubStation_Alpha
http://www.matroska.org/technical/specs/subtitles/ssa.html
http://moodub.free.fr/video/ass-specs.doc

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
	"strconv"
	"strings"
	"time"
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

// Info that determines a style. Contains the final format that goes into the SSA file.
type styleKey struct {
	Pos   int // Subtitle position
	Color int // Subtitle color, in BBGGRR format
}

// keyFromSub creates a style key containing the style info of the subtitle.
func keyFromSub(s *Subtitle) (k styleKey) {
	pos := s.Pos
	if pos == PosNotSpecified {
		pos = Bottom // Assign default position
	}
	k.Pos = modelPosToSsaPos[pos]

	color := s.Color

	for {
		if color == "" {
			// Unknown / unspecified color, assign a default value
			k.Color = 0xefefef // light gray
			break
		}
		if color[0] == '#' {
			color = color[1:]
		}
		// Try hex form
		if n, err := strconv.ParseInt(color, 16, 64); err == nil {
			// It's a hex form (RRGGBB). Switch bytes.
			v := int(n)
			k.Color = (v & 0xff0000) >> 16
			k.Color |= v & 0x00ff00
			k.Color |= v & 0x0000ff << 16
			break
		}
		// Not a hex form, try the standard color names
		color = htmlColorRGB[strings.ToLower(color)] // If unknown, will be "" and handled in next iteration
	}

	return
}

// WriteSsaTo generates Sub Station Alpha format and writes it to an io.Writer.
func WriteSsaTo(w io.Writer, sp *SubsPack) (err error) {
	wr := &writer{w: w}

	// Script Info section
	wr.prn("[Script Info]") // This must be the first line
	wr.prn("; This is a Sub Station Alpha v4 script.")
	wr.prn("; ", HomePage)
	wr.prn("Title: ")
	wr.prn("Script Updated By: srtgears version ", Version)
	wr.prn("ScriptType: v4.00")
	wr.prn("Collisions: Normal")
	wr.prn("PlayResY: 600")
	wr.prn("PlayDepth: 0")
	wr.prn("Timer: 100,0000")

	// Styles section
	wr.prn()
	wr.prn("[V4 Styles]")
	wr.prn("Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, TertiaryColour, BackColour, Bold, Italic, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, AlphaLevel, Encoding")
	// Loop over all subtitles to determine what styles we have
	styleKeys := make([]styleKey, len(sp.Subs)) // Store calculated style keys, we will need to interate over subs once more
	stylesMap := map[styleKey]string{}          // Unique style keys
	styles := []styleKey{}                      // Maintain order of unique style keys for generation
	for i, s := range sp.Subs {
		styleKeys[i] = keyFromSub(s)
		if _, ok := stylesMap[styleKeys[i]]; !ok {
			// new style
			styles = append(styles, styleKeys[i])
			stylesMap[styleKeys[i]] = strconv.Itoa(len(styles))
		}
	}
	// Now generate style definitions
	for _, v := range styles {
		wr.prf("Style: %s, Arial,28,%d,%d,%d,-2147483640,-1,0,1,1,2,%d,30,30,30,0,0",
			stylesMap[v], v.Color, v.Color, v.Color, v.Pos)
		wr.prn()
	}

	// Events section
	wr.prn()
	wr.prn("[Events]")
	wr.prn("Format: Marked, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text")

	printTime := func(t time.Duration) {
		hour := t / time.Hour
		min := (t % time.Hour) / time.Minute
		sec := (t % time.Minute) / time.Second
		ms := (t % time.Second) / time.Millisecond
		wr.prf("%d:%02d:%02d.%02d", hour, min, sec, ms/10)
	}

	for i, s := range sp.Subs {
		if wr.err != nil {
			break
		}

		wr.pr("Dialogue: Marked=0,")
		printTime(s.TimeIn)
		wr.pr(",")
		printTime(s.TimeOut)
		wr.pr(",", stylesMap[styleKeys[i]], ",NA,0000,0000,0000,,")

		// Texts
		for i, line := range s.Lines {
			// Note: HTML and controls not need to be removed (they will by the player)
			wr.pr(line)
			if i != len(s.Lines)-1 {
				wr.pr(`\n`)
			}
		}
		wr.prn()
	}

	return wr.err
}

// HTML color names and their RGB codes. SSA does not support color names.
// Src: http://www.w3schools.com/html/html_colornames.asp
var htmlColorRGB = map[string]string{
	"aliceblue":            "F0F8FF",
	"antiquewhite":         "FAEBD7",
	"aqua":                 "00FFFF",
	"aquamarine":           "7FFFD4",
	"azure":                "F0FFFF",
	"beige":                "F5F5DC",
	"bisque":               "FFE4C4",
	"black":                "000000",
	"blanchedalmond":       "FFEBCD",
	"blue":                 "0000FF",
	"blueviolet":           "8A2BE2",
	"brown":                "A52A2A",
	"burlywood":            "DEB887",
	"cadetblue":            "5F9EA0",
	"chartreuse":           "7FFF00",
	"chocolate":            "D2691E",
	"coral":                "FF7F50",
	"cornflowerblue":       "6495ED",
	"cornsilk":             "FFF8DC",
	"crimson":              "DC143C",
	"cyan":                 "00FFFF",
	"darkblue":             "00008B",
	"darkcyan":             "008B8B",
	"darkgoldenrod":        "B8860B",
	"darkgray":             "A9A9A9",
	"darkgrey":             "A9A9A9",
	"darkgreen":            "006400",
	"darkkhaki":            "BDB76B",
	"darkmagenta":          "8B008B",
	"darkolivegreen":       "556B2F",
	"darkorange":           "FF8C00",
	"darkorchid":           "9932CC",
	"darkred":              "8B0000",
	"darksalmon":           "E9967A",
	"darkseagreen":         "8FBC8F",
	"darkslateblue":        "483D8B",
	"darkslategray":        "2F4F4F",
	"darkslategrey":        "2F4F4F",
	"darkturquoise":        "00CED1",
	"darkviolet":           "9400D3",
	"deeppink":             "FF1493",
	"deepskyblue":          "00BFFF",
	"dimgray":              "696969",
	"dimgrey":              "696969",
	"dodgerblue":           "1E90FF",
	"firebrick":            "B22222",
	"floralwhite":          "FFFAF0",
	"forestgreen":          "228B22",
	"fuchsia":              "FF00FF",
	"gainsboro":            "DCDCDC",
	"ghostwhite":           "F8F8FF",
	"gold":                 "FFD700",
	"goldenrod":            "DAA520",
	"gray":                 "808080",
	"grey":                 "808080",
	"green":                "008000",
	"greenyellow":          "ADFF2F",
	"honeydew":             "F0FFF0",
	"hotpink":              "FF69B4",
	"indianred":            "CD5C5C",
	"indigo":               "4B0082",
	"ivory":                "FFFFF0",
	"khaki":                "F0E68C",
	"lavender":             "E6E6FA",
	"lavenderblush":        "FFF0F5",
	"lawngreen":            "7CFC00",
	"lemonchiffon":         "FFFACD",
	"lightblue":            "ADD8E6",
	"lightcoral":           "F08080",
	"lightcyan":            "E0FFFF",
	"lightgoldenrodyellow": "FAFAD2",
	"lightgray":            "D3D3D3",
	"lightgrey":            "D3D3D3",
	"lightgreen":           "90EE90",
	"lightpink":            "FFB6C1",
	"lightsalmon":          "FFA07A",
	"lightseagreen":        "20B2AA",
	"lightskyBlue":         "87CEFA",
	"lightslategray":       "778899",
	"lightslategrey":       "778899",
	"lightsteelBlue":       "B0C4DE",
	"lightyellow":          "FFFFE0",
	"lime":                 "00FF00",
	"limegreen":            "32CD32",
	"linen":                "FAF0E6",
	"magenta":              "FF00FF",
	"maroon":               "800000",
	"mediumaquamarine":     "66CDAA",
	"mediumblue":           "0000CD",
	"mediumorchid":         "BA55D3",
	"mediumpurple":         "9370DB",
	"mediumseagreen":       "3CB371",
	"mediumslateblue":      "7B68EE",
	"mediumspringgreen":    "00FA9A",
	"mediumturquoise":      "48D1CC",
	"mediumvioletred":      "C71585",
	"midnightblue":         "191970",
	"mintcream":            "F5FFFA",
	"mistyrose":            "FFE4E1",
	"moccasin":             "FFE4B5",
	"navajowhite":          "FFDEAD",
	"navy":                 "000080",
	"oldlace":              "FDF5E6",
	"olive":                "808000",
	"olivedrab":            "6B8E23",
	"orange":               "FFA500",
	"orangered":            "FF4500",
	"orchid":               "DA70D6",
	"palegoldenrod":        "EEE8AA",
	"palegreen":            "98FB98",
	"paleturquoise":        "AFEEEE",
	"palevioletred":        "DB7093",
	"papayawhip":           "FFEFD5",
	"peachpuff":            "FFDAB9",
	"peru":                 "CD853F",
	"pink":                 "FFC0CB",
	"plum":                 "DDA0DD",
	"powderblue":           "B0E0E6",
	"purple":               "800080",
	"rebeccapurple":        "663399",
	"red":                  "FF0000",
	"rosybrown":            "BC8F8F",
	"royalblue":            "4169E1",
	"saddlebrown":          "8B4513",
	"salmon":               "FA8072",
	"sandybrown":           "F4A460",
	"seagreen":             "2E8B57",
	"seashell":             "FFF5EE",
	"sienna":               "A0522D",
	"silver":               "C0C0C0",
	"skyblue":              "87CEEB",
	"slateblue":            "6A5ACD",
	"slategray":            "708090",
	"slategrey":            "708090",
	"snow":                 "FFFAFA",
	"springgreen":          "00FF7F",
	"steelblue":            "4682B4",
	"tan":                  "D2B48C",
	"teal":                 "008080",
	"thistle":              "D8BFD8",
	"tomato":               "FF6347",
	"turquoise":            "40E0D0",
	"violet":               "EE82EE",
	"wheat":                "F5DEB3",
	"white":                "FFFFFF",
	"whitesmoke":           "F5F5F5",
	"yellow":               "FFFF00",
	"yellowgreen":          "9ACD32",
}
