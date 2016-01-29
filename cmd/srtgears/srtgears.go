/*

This is the main package of the Srtgears command line tool.

*/
package main

import (
	"fmt"
	"github.com/icza/srtgears"
	"github.com/icza/srtgears/exec"
	"os"
	"path"
	"strings"
)

var Version = "dev" // Srtgears version, filled by build

// Our heart: the Executor
var e = exec.New(os.Stdout)

func main() {
	fmt.Printf("Srtgears %s, home page: %s\n", Version, srtgears.HomePage)

	e.FlagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of srtgears:\n")
		e.FlagSet.PrintDefaults()
		fmt.Println(examples)
	}

	if err := e.ProcFlags(os.Args[1:]); err != nil {
		return
	}

	if len(os.Args) < 2 {
		e.FlagSet.Usage()
	}

	if err := readFiles(); err != nil {
		fmt.Println(err)
		return
	}

	if err := e.GearIt(); err != nil {
		fmt.Println(err)
		return
	}

	if e.Stats {
		return // Stats modifies the subtitles, omit writing to files.
	}
	if err := writeFiles(); err != nil {
		fmt.Println(err)
		return
	}
}

// readFiles loads the subtitle files specified by the '-in' and '-in2' flags.
func readFiles() (err error) {
	if e.In != "" {
		if e.Sp1, err = srtgears.ReadSrtFile(e.In); err != nil {
			return
		}
	}
	if e.In2 != "" {
		if e.Sp2, err = srtgears.ReadSrtFile(e.In2); err != nil {
			return
		}
	}
	return
}

// writeFiles writes the output files specified by the '-out' and '-out2' flags.
func writeFiles() (err error) {
	wf := func(name string, sp *srtgears.SubsPack) (err error) {
		ext := strings.ToLower(path.Ext(name))
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

	if e.Out != "" && e.Sp1 != nil {
		if err = wf(e.Out, e.Sp1); err != nil {
			return
		}
	}
	if e.Out2 != "" && e.Sp2 != nil {
		if err = wf(e.Out2, e.Sp2); err != nil {
			return
		}
	}
	return
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
