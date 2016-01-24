/*

Contains configuration used by package srtgears.

*/

package srtgears

import (
	"log"
)

const (
	Version  = "1.0"                            // Srtgears engine version
	HomePage = "https://srt-gears.appspot.com/" // Srtgears home page
	Author   = "Andras Belicza"                 // Author name
)

var (
	Debug bool // Tells whether to print debug messages.
)

func debugf(format string, a ...interface{}) {
	if Debug {
		log.Printf("[DEBUG] "+format, a...)
	}
}
