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
	"net/http"
)

func init() {
	http.HandleFunc("/srtgears-online-submit", sgwHandler)

}

func sgwHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	c.Debugf("Location: %s;%s;%s;%s", r.Header.Get("X-AppEngine-Country"), r.Header.Get("X-AppEngine-Region"), r.Header.Get("X-AppEngine-City"), r.Header.Get("X-AppEngine-CityLatLong"))

	in, inh, err := r.FormFile("in")
	if err != nil {
		c.Errorf("No input srt file 'in': %v", err)
		fmt.Fprint(w, "You must select an input srt file!")
		return
	}
	c.Debugf("Received uploaded file 'in': %s", inh.Filename)

	in2, inh2, err := r.FormFile("in2")
	if err != nil {
		c.Debugf("No 2nd input srt 'in2': %v", err)
	} else {
		c.Debugf("Received uploaded file 'in2': %s", inh2.Filename)
	}

	var sp1, sp2 *srtgears.SubsPack

	if sp1, err = srtgears.ReadSrtFrom(in); err != nil {
		c.Errorf("Failed to parse uploaded file 'in': %v", err)
		fmt.Fprint(w, "Failed to parse uploaded file: ", err)
		return
	}

	if in2 != nil {
		if sp2, err = srtgears.ReadSrtFrom(in2); err != nil {
			c.Errorf("Failed to parse uploaded file 'in2': %v", err)
			fmt.Fprint(w, "Failed to parse 2nd uploaded file: ", err)
			return
		}
	}

	_, _ = sp1, sp2

	fmt.Fprint(w, "Form submission received.")
}
