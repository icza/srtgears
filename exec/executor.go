/*

Package exec is a parameter driven executor which parses the parameters and performs subtitle transformations.

It is the heart of both the the srtgears command line tool and the web interface
(the web interface produces an argument list you would provide in the command line).

*/
package exec

import (
	"flag"
	"fmt"
	"github.com/gophergala2016/srtgears"
	"regexp"
	"strconv"
	"time"
)

// Executor is a helper type of which you can execute subtitle transformations defined by a series of arguments.
type Executor struct {
	FlagSet *flag.FlagSet // Custom Flagset used to parse parameters

	In         string  // input file name (*.srt)
	Out        string  // output file name (*.srt or *.ssa)
	In2        string  // optional 2nd input file name (when merging or concatenating subtitles) (*.srt)
	Out2       string  // optional 2nd output file name (when splitting) (*.srt or *.ssa)
	Concat     string  // concatenate 2 subtitle files, 2nd part start at e.g. '00:59:00,123'
	Merge      bool    // merge 2 subtitle files ('-in' at bottom, '-in2' at top
	SplitAt    string  // time where to split to 2 subtitle files ('-out' and '-out2'), e.g. '00:59:00,123'
	ShiftBy    int     // shift subtitle timestamps (+/- ms)
	Scale      float64 // scale subtitle timestamps (faster/slower); multiplier e.g. 1.001
	Lengthen   float64 // lengthen / shorten display duration of subtitles, multiplier e.g. for +10% use 1.1
	RemoveHTML bool    // strip off formatting (e.g. <i>, <b>, <u>, <font> etc.)
	RemoveCtrl bool    // remove controls such as {\anX} (or {\aY}), {\pos(x,y)}
	RemoveHI   bool    // remove hearing impaired subtitles (such as '[PHONE RINGING]' or '(phone ringing)')
	Pos        string  // change subtitle position, one of: BL, B, BR, L, C, R, TL, T, TR  (B: bottom, T: Top, L: Left, R: Right, C: Center)
	Color      string  // change subtitle color, name (e.g. 'red' or 'yellow') or RGB hexa '#rrggbb' (e.g.'#ff0000' for red)
	Stats      bool    // analyze file and print statistics

	Sp1, Sp2 *srtgears.SubsPack // SubsPacks to operate on. Must be set by the user before calling GearIt()!
}

// New creates a new Executor.
func New() *Executor {
	return &Executor{
		FlagSet: flag.NewFlagSet("name", flag.ContinueOnError),
	}
}

// ProcFlags sets up variables for parsing the arguments, pointing to the fields of the Executor.
// And parses the arguments.
func (e *Executor) ProcFlags(arguments []string) error {
	f := e.FlagSet

	f.StringVar(&e.In, "in", "", "input file name (*.srt)")
	f.StringVar(&e.Out, "out", "", "output file name (*.srt or *.ssa)")
	f.StringVar(&e.In2, "in2", "", "optional 2nd input file name (when merging or concatenating subtitles) (*.srt)")
	f.StringVar(&e.Out2, "out2", "", "optional 2nd output file name (when splitting) (*.srt or *.ssa)")
	f.BoolVar(&srtgears.Debug, "debug", false, "print debug messages")
	f.StringVar(&e.Concat, "concat", "", "concatenate 2 subtitle files, 2nd part start at e.g. '00:59:00,123'")
	f.BoolVar(&e.Merge, "merge", false, "merge 2 subtitle files ('-in' at bottom, '-in2' at top)")
	f.StringVar(&e.SplitAt, "splitAt", "", "time where to split to 2 subtitle files ('-out' and '-out2'), e.g. '00:59:00,123'")
	f.IntVar(&e.ShiftBy, "shiftBy", 0, "shift subtitle timestamps (+/- ms)")
	f.Float64Var(&e.Scale, "scale", 0, "scale subtitle timestamps (faster/slower); multiplier e.g. 1.001")
	f.Float64Var(&e.Lengthen, "lengthen", 0, "lengthen / shorten display duration of subtitles, multiplier e.g. for +10% use 1.1")
	f.BoolVar(&e.RemoveHTML, "removehtml", false, "strip off formatting (e.g. <i>, <b>, <u>, <font> etc.)")
	f.BoolVar(&e.RemoveCtrl, "removectrl", false, `remove controls such as {\anX} (or {\aY}), {\pos(x,y)}`)
	f.BoolVar(&e.RemoveHI, "removehi", false, "remove hearing impaired subtitles (such as '[PHONE RINGING]' or '(phone ringing)')")
	f.StringVar(&e.Pos, "pos", "", "change subtitle position, one of: BL, B, BR, L, C, R, TL, T, TR  (B: bottom, T: Top, L: Left, R: Right, C: Center)")
	f.StringVar(&e.Color, "color", "", "change subtitle color, name (e.g. 'red' or 'yellow') or RGB hexa '#rrggbb' (e.g.'#ff0000' for red)")
	f.BoolVar(&e.Stats, "stats", false, "analyze file and print statistics")

	return f.Parse(arguments)
}

// Regexp pattern used to parse timestamps.
var timestampPattern = regexp.MustCompile(`(\d\d):(\d\d):(\d\d)[,\.](\d\d\d)`)

// parseTime parses a timestamp which must be in the form of
// 00:00:00,000
func parseTime(t string) (time.Duration, error) {
	parts := timestampPattern.FindStringSubmatch(t)
	if len(parts) == 0 {
		return 0, fmt.Errorf("Invalid time: %s", t)
	}
	get := func(idx int) time.Duration {
		n, err := strconv.ParseInt(parts[idx], 10, 64)
		if err != nil {
			panic(err) // This shouldn't happen as only digits are matched.
		}
		return time.Duration(n)
	}
	return time.Hour*get(1) + time.Minute*get(2) + time.Second*get(3) + time.Millisecond*get(4), nil
}

// Mapping between positions expected in arguments to our model Pos.
var argPosToModelPos = map[string]srtgears.Pos{
	"TL": srtgears.TopLeft, "T": srtgears.Top, "TR": srtgears.TopRight,
	"L": srtgears.Left, "C": srtgears.Center, "R": srtgears.Right,
	"BL": srtgears.BottomLeft, "B": srtgears.Bottom, "BR": srtgears.BottomRight,
}

// GearIt performs subtitle transformations specified by the arguments passed to ProcFlags().
// Prior to calling this method, Executor.Sp1 and Executor.Sp2 should be set.
func (e *Executor) GearIt() (err error) {
	sp1, sp2 := e.Sp1, e.Sp2

	if sp1 == nil {
		return fmt.Errorf("Input file must be specified ('-in')!")
	}
	if sp2 == nil && (e.Concat != "" || e.Merge) {
		return fmt.Errorf("2nd input file must be specified ('-in2')!")
	}

	if e.Concat != "" {
		secPartStart, err := parseTime(e.Concat)
		if err != nil {
			return fmt.Errorf("Invalid time for concat: %s", e.Concat)
		}
		sp1.Concatenate(sp2, secPartStart)
	}

	if e.Merge {
		sp1.Merge(sp2)
	}

	if e.Lengthen != 0 {
		sp1.Lengthen(e.Lengthen)
	}

	if e.RemoveCtrl {
		sp1.RemoveControl()
	}

	if e.RemoveHI {
		sp1.RemoveHI()
	}

	if e.RemoveHTML {
		sp1.RemoveHTML()
	}

	if e.Pos != "" {
		pos2, ok := argPosToModelPos[e.Pos]
		if !ok {
			return fmt.Errorf("Invalid pos value: %s", e.Pos)
		}
		sp1.SetPos(pos2)
	}

	if e.Color != "" {
		sp1.SetColor(e.Color)
	}

	if e.Scale != 0 {
		sp1.Scale(e.Scale)
	}

	if e.ShiftBy != 0 {
		sp1.Shift(time.Duration(e.ShiftBy) * time.Millisecond)
	}

	if e.SplitAt != "" {
		at, err := parseTime(e.SplitAt)
		if err != nil {
			return fmt.Errorf("Invalid time for splitAt: %s", e.SplitAt)
		}
		sp2 = sp1.Split(at)
		e.Sp2 = sp2 // sp2 is just a local copy, so we need to update Executor.Sp2 too!
	}

	if e.Stats {
		ss := sp1.Stats()
		fmt.Printf("STATS of %s:\n", e.In)
		p := func(name string, value interface{}) {
			fmt.Printf("%-29s: %v\n", name, value)
		}
		p("Total # of subtitles", ss.Subs)
		p("Lines", ss.Lines)
		p("Avg lines per sub", fmt.Sprintf("%.4f", ss.AvgLinesPerSub))
		p("Chars (with spaces)", ss.Chars)
		p("Chars (without spaces)", ss.CharsNoSpace)
		p("Avg chars (no space) per line", fmt.Sprintf("%.4f", ss.AvgCharsPerLine))
		p("Words", ss.Words)
		p("Avg words per line", fmt.Sprintf("%.4f", ss.AvgWordsPerLine))
		p("Avg chars per word", fmt.Sprintf("%.4f", ss.AvgCharsPerWord))
		p("Total subtitle display time", ss.TotalDispDur)
		p("Subtitle visible ratio", fmt.Sprintf("%.2f%% (compared to total length)", ss.SubVisibRatio*100))
		p("Avg. display duration", ss.AvgDispDurPerNonSpaceChar.String()+" per 1 non-space char")
		p("Subs with HTML formatting", ss.HTMLs)
		p("Subs with controls", ss.Controls)
		p("Subs with hearing impaired", ss.HIs)
	}

	return
}
