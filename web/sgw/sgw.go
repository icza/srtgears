/*

Package sgw is the backend logic for the online Srtgears.

It processes form submit requests, and connects the data to the srtgears engine,
and sends back the result.

*/
package sgw

import (
	"appengine"
	"archive/zip"
	"fmt"
	"github.com/gophergala2016/srtgears"
	"github.com/gophergala2016/srtgears/exec"
	"io"
	"net/http"
	"path"
	"strings"
	"time"
)

func init() {
	http.HandleFunc("/srtgears-online-submit", sgwHandler)
}

// sgwHandler handles the form submits.
func sgwHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	c.Debugf("Location: %s;%s;%s;%s", r.Header.Get("X-AppEngine-Country"), r.Header.Get("X-AppEngine-Region"),
		r.Header.Get("X-AppEngine-City"), r.Header.Get("X-AppEngine-CityLatLong"))

	args := []string{} // To simulate command line arguments

	// Form rewind: fill args slice with the posted form values

	in, inh, err := r.FormFile("in")
	if err != nil {
		c.Errorf("No input srt file 'in': %v", err)
		fmt.Fprint(w, "You must select an input srt file!")
		return
	}
	c.Debugf("Received uploaded file 'in': %s", inh.Filename)
	args = append(args, "-in", inh.Filename)

	in2, inh2, err := r.FormFile("in2")
	if err != nil {
		c.Debugf("No 2nd input srt 'in2': %v", err)
	} else {
		c.Debugf("Received uploaded file 'in2': %s", inh2.Filename)
		args = append(args, "-in2", inh2.Filename)
	}

	args = rewindForm(args, r)

	// Our heart: the Executor
	var e = exec.New(w)         // If there are errors, we want them generated on the response.
	e.FlagSet.Usage = func() {} // We don't want usage in the web response

	if err := e.ProcFlags(args); err != nil {
		return // Errors are already written to response.
	}

	// Read input files
	if e.Sp1, err = srtgears.ReadSrtFrom(in); err != nil {
		c.Errorf("Failed to parse uploaded file 'in': %v", err)
		fmt.Fprint(w, "Failed to parse uploaded file: ", err)
		return
	}

	if in2 != nil {
		if e.Sp2, err = srtgears.ReadSrtFrom(in2); err != nil {
			c.Errorf("Failed to parse uploaded file 'in2': %v", err)
			fmt.Fprint(w, "Failed to parse 2nd uploaded file: ", err)
			return
		}
	}

	// We want stats in plain text...
	e.BeforeStats = func() {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	// Perform transformations
	if err := e.GearIt(); err != nil {
		fmt.Fprint(w, err)
		return
	}

	if e.Stats {
		return // If stats was specified, response is already committed.
	}

	// Everything went ok. We can now generate and send the transformed subtitles.
	if err := sendSubs(w, e); err != nil {
		c.Errorf("Failed to send subtitles: %v", err)
	}
}

// sendSubs generates and send the transformed subtitles, zipped.
func sendSubs(w http.ResponseWriter, e *exec.Executor) (err error) {
	// First checks extensions so we can send back error.
	// Once we start writing zip, there's no going back.
	validExt := func(name string) bool {
		switch ext := strings.ToLower(path.Ext(name)); ext {
		case ".srt", ".ssa":
			return true
		case "":
			fmt.Fprintf(w, "Output extension not specified: %s", name)
		default:
			fmt.Fprintf(w, "Unsupported file extension, only *.srt and *.ssa are supported: %s", ext)
		}
		return false
	}

	fileCount := 0
	if e.Out != "" && e.Sp1 != nil {
		if !validExt(e.Out) {
			return
		}
		fileCount++
	}
	if e.Out2 != "" && e.Sp2 != nil {
		if !validExt(e.Out2) {
			return
		}
		fileCount++
	}

	if fileCount == 0 { // Just to make sure
		fmt.Fprint(w, "No output file has been specified.")
		return
	}
	if fileCount == 2 && e.Out == e.Out2 {
		fmt.Fprint(w, "The 2 output file names cannot be the same!")
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=subpack.zip")

	zw := zip.NewWriter(w)
	defer zw.Close()

	wf := func(name string, sp *srtgears.SubsPack) (err error) {
		var f io.Writer
		fh := &zip.FileHeader{Name: name}
		fh.SetModTime(time.Now())
		if f, err = zw.CreateHeader(fh); err != nil {
			return
		}
		switch ext := strings.ToLower(path.Ext(name)); ext {
		case ".srt":
			return srtgears.WriteSrtTo(f, sp)
		case ".ssa":
			return srtgears.WriteSsaTo(f, sp)
		}
		return
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

// rewindForm reads form values and generates proper arguments for them.
// Returns the updated args slice.
func rewindForm(args []string, r *http.Request) []string {
	if s := r.FormValue("out"); s != "" {
		args = append(args, "-out", s)
	}
	if s := r.FormValue("out2"); s != "" {
		args = append(args, "-out2", s)
	}
	if s := r.FormValue("concat"); s != "" {
		args = append(args, "-concat="+s)
	}
	if s := r.FormValue("merge"); s != "" {
		args = append(args, "-merge")
	}
	if s := r.FormValue("lengthen"); s != "" {
		args = append(args, "-lengthen="+s)
	}
	if s := r.FormValue("removectrl"); s != "" {
		args = append(args, "-removectrl")
	}
	if s := r.FormValue("removehi"); s != "" {
		args = append(args, "-removehi")
	}
	if s := r.FormValue("removehtml"); s != "" {
		args = append(args, "-removehtml")
	}
	if s := r.FormValue("pos"); s != "" {
		args = append(args, "-pos="+s)
	}
	if s := r.FormValue("color"); s != "" {
		args = append(args, "-color="+s)
	}
	if s := r.FormValue("scale"); s != "" {
		args = append(args, "-scale="+s)
	}
	if s := r.FormValue("shiftBy"); s != "" {
		args = append(args, "-shiftBy="+s)
	}
	if s := r.FormValue("splitAt"); s != "" {
		args = append(args, "-splitAt="+s)
	}
	if s := r.FormValue("stats"); s != "" {
		args = append(args, "-stats")
	}

	return args
}
