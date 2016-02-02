/*

This is the main package of the packrelease tool which packs cross-compiled releases
into the web/release folder, and also generates the HTML table in the format ready for the
download.html.

*/
package main

import (
	"log"
	"os"
	"path/filepath"
)

var root = "../srtgears"

func main() {
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	root, err := filepath.Abs(root)
	if err != nil {
		panic(err)
	}
	log.Println("Using root:", root)
	
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		log.Println("Found:", path)
		return nil
	})
}
