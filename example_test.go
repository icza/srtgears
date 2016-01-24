/*

Example functions.

*/

package srtgears

// Example shows how to merge 2 subtitle files to have a dual sub with srtgears.
func Example() {
	sp1, err := srtgears.ReadSrtFile("eng.srt")
	// check err
	sp2, err := srtgears.ReadSrtFile("hun.srt")
	// check err
	sp1.Merge(sp2)
	err = srtgears.WriteSrtFile("eng+hun.srt", sp1)
	// check err
}
