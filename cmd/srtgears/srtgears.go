/*

This is the main package of the Srtgears command line tool.

*/
package main

import (
	"flag"
	"fmt"
	"github.com/gophergala2016/srtgears"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	in         string  // input file name (*.srt)
	out        string  // output file name (*.srt or *.ssa)
	in2        string  // optional 2nd input file name (when merging or concatenating subtitles) (*.srt)
	out2       string  // optional 2nd output file name (when splitting) (*.srt or *.ssa)
	concat     string  // concatenate 2 subtitle files, 2nd part start at e.g. '00:59:00,123'
	merge      bool    // merge 2 subtitle files ('-in' at bottom, '-in2' at top
	splitAt    string  // time where to split to 2 subtitle files ('-out' and '-out2'), e.g. '00:59:00,123'
	shiftBy    int     // shift subtitle timestamps (+/- ms)
	scale      float64 // scale subtitle timestamps (faster/slower); multipler e.g. 1.001
	lengthen   float64 // lengthen / shorten display duration of subtitles, multipler e.g. for +10% use 1.1
	removeHTML bool    // strip off formatting (e.g. <i>, <b>, <u>, <font> etc.)
	removeCtrl bool    // remove controls such as {\anX} (or {\aY}), {\pos(x,y)}
	removeHI   bool    // remove hearing impaired subtitles (such as '[PHONE RINGING]' or '(phone ringing)')
	pos        string  // change subtitle position, one of: BL, B, BR, L, C, R, TL, T, TR  (B: bottom, T: Top, L: Left, R: Right, C: Center)
	color      string  // change subtitle color, name (e.g. 'red' or 'yellow') or RGB hexa '#rrggbb' (e.g.'#ff0000' for red)
	stats      bool    // analyze file and print statistics
)

// Loaded subtitles
var sp1, sp2 *srtgears.SubsPack

func main() {
	fmt.Printf("Srtgears %s, home page: %s\n", srtgears.Version, srtgears.HomePage)

	if err := procFlags(); err != nil {
		fmt.Println(err)
		flag.Usage()
		return
	}

	if err := readFiles(); err != nil {
		fmt.Println(err)
		return
	}

	if err := gearIt(sp1, sp2); err != nil {
		fmt.Println(err)
		return
	}

	if err := writeFiles(); err != nil {
		fmt.Println(err)
		return
	}
}

// readFiles loads the subtitle files specified by the '-in' and '-in2' flags.
func readFiles() (err error) {
	if in != "" {
		if sp1, err = srtgears.ReadSrtFile(in); err != nil {
			return
		}
	}
	if in2 != "" {
		if sp2, err = srtgears.ReadSrtFile(in2); err != nil {
			return
		}
	}
	return
}

// parseTime parses a timestamp which must be in the form of
// 00:00:00,000
func parseTime(t string) (time.Duration, error) {
	parts := regexp.MustCompile(`(\d\d):(\d\d):(\d\d)[,\.](\d\d\d)`).FindStringSubmatch(t)
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

// gearIt performs subtitle transformations specified by the command line flags
func gearIt(sp1, sp2 *srtgears.SubsPack) (err error) {
	if sp1 == nil {
		return fmt.Errorf("Input file must be specified ('-in')!")
	}
	if sp2 == nil && (concat != "" || merge) {
		return fmt.Errorf("2nd input file must be specified ('-in2')!")
	}

	if concat != "" {
		secPartStart, err := parseTime(concat)
		if err != nil {
			return fmt.Errorf("Invalid time for concat: %s", concat)
		}
		sp1.Concatenate(sp2, secPartStart)
	}

	if merge {
		sp1.Merge(sp2)
	}

	if lengthen != 0 {
		sp1.Lengthen(lengthen)
	}

	if removeCtrl {
		sp1.RemoveControl()
	}

	if removeHI {
		sp1.RemoveHI()
	}

	if removeHTML {
		sp1.RemoveHTML()
	}

	if pos != "" {
		m := map[string]srtgears.Pos{
			"TL": srtgears.TopLeft, "T": srtgears.Top, "TR": srtgears.TopRight,
			"L": srtgears.Left, "C": srtgears.Center, "R": srtgears.Right,
			"BL": srtgears.BottomLeft, "B": srtgears.Bottom, "BR": srtgears.BottomRight,
		}
		pos2, ok := m[pos]
		if !ok {
			return fmt.Errorf("Invalid pos value: %s", pos)
		}
		sp1.SetPos(pos2)
	}

	if color != "" {
		sp1.SetColor(color)
	}

	if scale != 0 {
		sp1.Scale(scale)
	}

	if shiftBy != 0 {
		sp1.Shift(time.Duration(shiftBy) * time.Millisecond)
	}

	if splitAt != "" {
		at, err := parseTime(concat)
		if err != nil {
			return fmt.Errorf("Invalid time for splitAt: %s", splitAt)
		}
		sp2 = sp1.Split(at)
	}

	if stats {
		ss := sp1.Stats()
		fmt.Printf("STATS of %s:\n", in)
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

// writeFiles writes the output files specified by the '-out' and '-out2' flags.
func writeFiles() (err error) {
	wf := func(name string, sp *srtgears.SubsPack) (err error) {
		ext := strings.ToLower(path.Ext(out))
		switch ext {
		case ".srt":
			return srtgears.WriteSrtFile(name, sp)
		case ".ssa":
			return srtgears.WriteSsaFile(name, sp)
		case "":
			return fmt.Errorf("Output extension not specified!")
		}
		return fmt.Errorf("Unsupported file extension, only *.srt and *.ssa are supported: %s", ext)
	}

	if out != "" && sp1 != nil {
		if err = wf(out, sp1); err != nil {
			return
		}
	}
	if out2 != "" && sp2 != nil {
		if err = wf(out2, sp2); err != nil {
			return
		}
	}
	return
}

func procFlags() error {
	flag.StringVar(&in, "in", "", "input file name (*.srt)")
	flag.StringVar(&out, "out", "", "output file name (*.srt or *.ssa)")
	flag.StringVar(&in2, "in2", "", "optional 2nd input file name (when merging or concatenating subtitles) (*.srt)")
	flag.StringVar(&out2, "out2", "", "optional 2nd output file name (when splitting) (*.srt or *.ssa)")
	flag.BoolVar(&srtgears.Debug, "debug", false, "print debug messages")
	flag.StringVar(&concat, "concat", "", "concatenate 2 subtitle files, 2nd part start at e.g. '00:59:00,123'")
	flag.BoolVar(&merge, "merge", false, "merge 2 subtitle files ('-in' at bottom, '-in2' at top")
	flag.StringVar(&splitAt, "splitAt", "", "time where to split to 2 subtitle files ('-out' and '-out2'), e.g. '00:59:00,123'")
	flag.IntVar(&shiftBy, "shiftBy", 0, "shift subtitle timestamps (+/- ms)")
	flag.Float64Var(&scale, "scale", 0, "scale subtitle timestamps (faster/slower); multipler e.g. 1.001")
	flag.Float64Var(&lengthen, "lengthen", 0, "lengthen / shorten display duration of subtitles, multipler e.g. for +10% use 1.1")
	flag.BoolVar(&removeHTML, "removehtml", false, "strip off formatting (e.g. <i>, <b>, <u>, <font> etc.)")
	flag.BoolVar(&removeCtrl, "removectrl", false, `remove controls such as {\anX} (or {\aY}), {\pos(x,y)}`)
	flag.BoolVar(&removeHI, "removehi", false, "remove hearing impaired subtitles (such as '[PHONE RINGING]' or '(phone ringing)')")
	flag.StringVar(&pos, "pos", "", "change subtitle position, one of: BL, B, BR, L, C, R, TL, T, TR  (B: bottom, T: Top, L: Left, R: Right, C: Center)")
	flag.StringVar(&color, "color", "", "change subtitle color, name (e.g. 'red' or 'yellow') or RGB hexa '#rrggbb' (e.g.'#ff0000' for red)")
	flag.BoolVar(&stats, "stats", false, "analyze file and print statistics")

	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		fmt.Println(examples)
	}

	flag.Parse()

	return nil
}

const examples = `
Examples:
Merge 2 files to have a dual sub saved in Sub Station Alpha (*.ssa) format:
    srtgears -in eng.srt -in2 hun.srt -out eng+hun.ssa
Concatenate 2 files where 2nd part of the movie starts at 51 min 15 sec:
    srtgears -in cd1.srt -in2 cd2.srt -out cd12.srt -concat=00:51:15:00,000
Change subtitle color to yellow, move to top, remove HI lines, increase display duration by 10% and save as *.ssa:
    srtgears -in eng.srt -out eng2.ssa -color=yellow -pos=T -removehi -lengthen=1.1
Repair: do nothing, just parse and re-save
    srtgears -in eng.srt -out eng2.srt`
