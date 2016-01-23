/*

This is the main package of the Srtgears command line tool.

*/
package main

import (
	"flag"
	"fmt"
	"github.com/gophergala2016/srtgears"
	"os"
)

const (
	version  = "1.0"                            // Srtgears application version
	homePage = "https://srt-gears.appspot.com/" // Srtgears home page
)

func main() {
	flag.Parse()
	fmt.Printf("Srtgears %s, home page: %s\n", version, homePage)

	sp, err := srtgears.ReadSrtFile("w:/video/a.srt")
	if err != nil {
		panic(err)
	}
	if err = srtgears.WriteSrtTo(os.Stdout, sp); err != nil {
		panic(err)
	}
}
