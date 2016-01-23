/*

This is the main package of the Srtgears command line tool.

*/
package main

import (
	"flag"
	"fmt"
	"github.com/gophergala2016/srtgears"
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
	lengthen   float64 // lengthen / shorten display duration of subtitles, multipler e.g. 1.1 for +10%
	removeHTML bool    // strip off formatting (e.g. <i>, <b>, <u>, <font> etc.)
	removeCtrl bool    // remove controls such as {\anX} (or {\aY}), {\pos(x,y)}
	removeHI   bool    // remove hearing impaired subtitles (such as '[PHONE RINGING]' or '(phone ringing)')
	pos        string  // change subtitle position, one of: BL, B, BR, L, C, R, TL, T, TR  (B: bottom, T: Top, L: Left, R: Right, C: Center)
	color      string  // change subtitle color, name (e.g. 'red' or 'yellow') or RGB hexa '#rrggbb' (e.g.'#ff0000' for red)
	repair     bool    // do nothing just parse and re-save
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

	if readFiles() != nil {
		return
	}

	gearIt(sp1, sp2)

	if writeFiles() != nil {
		return
	}
}

// readFiles loads the subtitle files specified by the '-in' and '-in2' flags.
func readFiles() (err error) {
	if in != "" {
		if sp1, err = srtgears.ReadSrtFile(in); err != nil {
			fmt.Println("Failed reading %s: %v", in, err)
			return
		}
	}
	if in2 != "" {
		if sp2, err = srtgears.ReadSrtFile(in2); err != nil {
			fmt.Println("Failed reading %s: %v", in2, err)
			return
		}
	}
	return
}

// gearIt performs subtitle transformations specified by the command line flags
func gearIt(sp1, sp2 *srtgears.SubsPack) {
	if merge {
		sp1.Merge(sp2)
	}
}

// writeFiles writes the output files specified by the '-out' and '-out2' flags.
func writeFiles() (err error) {
	if out != "" && sp1 != nil {
		if err = srtgears.WriteSrtFile(out, sp1); err != nil {
			fmt.Println("Failed writing %s: %v", out, err)
			return
		}
	}
	if out2 != "" && sp2 != nil {
		if err = srtgears.WriteSrtFile(out2, sp2); err != nil {
			fmt.Println("Failed writing %s: %v", out2, err)
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
	flag.StringVar(&concat, "concat", "00:00:00,000", "concatenate 2 subtitle files, 2nd part start at e.g. '00:59:00,123'")
	flag.BoolVar(&merge, "merge", false, "merge 2 subtitle files ('-in' at bottom, '-in2' at top")
	flag.StringVar(&splitAt, "splitAt", "00:00:00,000", "time where to split to 2 subtitle files ('-out' and '-out2'), e.g. '00:59:00,123'")
	flag.IntVar(&shiftBy, "shiftBy", 0, "shift subtitle timestamps (+/- ms)")
	flag.Float64Var(&scale, "scale", 0.0, "scale subtitle timestamps (faster/slower); multipler e.g. 1.001")
	flag.Float64Var(&lengthen, "lengthen", 0.0, "lengthen / shorten display duration of subtitles, multipler e.g. 1.1 for +10%")
	flag.BoolVar(&removeHTML, "removehtml", false, "strip off formatting (e.g. <i>, <b>, <u>, <font> etc.)")
	flag.BoolVar(&removeCtrl, "removectrl", false, `remove controls such as {\anX} (or {\aY}), {\pos(x,y)}`)
	flag.BoolVar(&removeHI, "removehi", false, "remove hearing impaired subtitles (such as '[PHONE RINGING]' or '(phone ringing)')")
	flag.StringVar(&pos, "pos", "", "change subtitle position, one of: BL, B, BR, L, C, R, TL, T, TR  (B: bottom, T: Top, L: Left, R: Right, C: Center)")
	flag.StringVar(&color, "color", "", "change subtitle color, name (e.g. 'red' or 'yellow') or RGB hexa '#rrggbb' (e.g.'#ff0000' for red)")
	flag.BoolVar(&repair, "repair", false, "do nothing just parse and re-save")
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
To merge 2 files to have a dual sub:
    srtgears -in eng.srt -in2 hun.srt -out eng+hun.srt
To change color to yellow, move subtitles to top, remove hearing impaired subtitles and increase display duration by 10%:
    srtgears -in eng.srt -out eng2.srt -color=yellow -pos=T -removehi -lengthen=1.1`
