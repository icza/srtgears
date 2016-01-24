/*

A simple writer-utility to ease error handling and specific newlines when
writing subtitle files.

*/

package srtgears

import (
	"fmt"
	"io"
)

const newLine = "\r\n" // Use Windows-style newline

// Writer to ease error handling and specific newline.
type writer struct {
	w   io.Writer // Destination writer
	err error     // Errors returned by w are stored here
}

// pr forwards to fmt.Fprint() if there were no errors before.
func (w *writer) pr(a ...interface{}) {
	if w.err == nil {
		_, w.err = fmt.Fprint(w.w, a...)
	}
}

// prn calls pr() and calls it again to print a newline.
func (w *writer) prn(a ...interface{}) {
	w.pr(a...)
	w.pr(newLine)
}

// prf forwards to fmt.Fprintf() if there were no errors before.
func (w *writer) prf(format string, a ...interface{}) {
	if w.err == nil {
		_, w.err = fmt.Fprintf(w.w, format, a...)
	}
}
