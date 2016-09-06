/*

Contains configuration used by package srtgears.

*/

package srtgears

import (
	"log"
)

const (
	// Srtgears home page
	HomePage = "https://srt-gears.appspot.com/"
	// Author name
	Author = "Andras Belicza"
)

var (
	// Tells whether to print debug messages.
	Debug bool
)

func debugf(format string, a ...interface{}) {
	if Debug {
		log.Printf("[DEBUG] "+format, a...)
	}
}
