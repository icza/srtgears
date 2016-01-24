/*

Package sgw is the backend logic for the online Srtgears.

It processes form submit requests, and connects the data to the srtgears engine,
and sends back the result.

*/
package sgw

import (
	"appengine"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/srtgears-online-submit", sgwHandler)

}

func sgwHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	c.Debugf("Location: %s;%s;%s;%s", r.Header.Get("X-AppEngine-Country"), r.Header.Get("X-AppEngine-Region"), r.Header.Get("X-AppEngine-City"), r.Header.Get("X-AppEngine-CityLatLong"))
	
	fmt.Fprint(w, "Form submit received.")
}
