/*

This is the main package of the Srtgears command line tool.

*/
package main

import (
	"fmt"
	"flag"
)

const version = "1.0" // Strgears application version

func main() {
	flag.Parse()
	fmt.Println("Strgears", version)
}
