/*

Package sgw is the backend logic for the online Srtgears.

It processes form submit requests, and connects the data to the srtgears engine,
and sends back the result.

*/
package sgw

import (
	"appengine"
	"fmt"
	"github.com/gophergala2016/srtgears"
	"github.com/gophergala2016/srtgears/exec"
	"net/http"
)

func init() {
	http.HandleFunc("/srtgears-online-submit", sgwHandler)
}

func sgwHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	c.Debugf("Location: %s;%s;%s;%s", r.Header.Get("X-AppEngine-Country"), r.Header.Get("X-AppEngine-Region"), r.Header.Get("X-AppEngine-City"), r.Header.Get("X-AppEngine-CityLatLong"))

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

	// Perform transformations
	if err := e.GearIt(); err != nil {
		fmt.Fprint(w, err)
		return
	}

	// Send result, zipped
	// TODO
}
