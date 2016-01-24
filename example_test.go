/*

Example functions.

*/

package srtgears

// This example shows how to merge 2 subtitle files to have a dual sub saved in Sub Station Alpha (*.ssa) format.
func Example_merge() {
	sp1, err := srtgears.ReadSrtFile("eng.srt")
	// check err
	sp2, err := srtgears.ReadSrtFile("hun.srt")
	// check err
	sp1.Merge(sp2)
	err = srtgears.WriteSsaFile("eng+hun.ssa", sp1)
	// check err
}

// This example shows how to concatenate 2 files where 2nd part of the movie starts at 51 min 15 sec.
func Example_concat() {
	sp1, err := srtgears.ReadSrtFile("cd1.srt")
	// check err
	sp2, err := srtgears.ReadSrtFile("cd2.srt")
	// check err
	secPartStart := time.Minute*51 + time.Second*15
	sp1.Concatenate(sp2, secPartStart)
	err = srtgears.WriteSrtFile("cd12.srt", sp1)
	// check err
}

// This example shows how to change subtitle color to yellow,
// move to top, remove HI (hearing impaired lines),
// increase display duration by 10% and save result as a
// Sub Station Alpha (*.ssa) file.
func Example_misc() {
	sp1, err := srtgears.ReadSrtFile("eng.srt")
	// check err
	sp1.SetColor("yellow")
	sp1.SetPos(srtgears.Top)
	sp1.RemoveHI()
	sp1.Lengthen(1.1)
	err = srtgears.WriteSsaFile("eng2.ssa", sp1)
	// check err
}
