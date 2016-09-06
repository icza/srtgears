/*

Contains configuration used by package srtgears.

*/

package srtgears

import (
	"log"
)

const (
	// HomePage is the Srtgears home page
	HomePage = "https://srt-gears.appspot.com/"
	// Author name
	Author = "Andras Belicza"
)

var (
	// Debug tells whether to print debug messages.
	Debug bool
)

func debugf(format string, a ...interface{}) {
	if Debug {
		log.Printf("[DEBUG] "+format, a...)
	}
}
