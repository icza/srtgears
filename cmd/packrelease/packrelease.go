/*

This is the main package of the packrelease tool which packs cross-compiled releases
into the web/release folder, and also generates the HTML table in the format ready for the
download.html.

*/
package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

// Release folder name pattern, e.g. "srtgears-1.1-windows-amd64"
var rlsFldrPttrn = regexp.MustCompile("srtgears-([^-]*)-([^-]*)-([^-]*)")

// Folders relative to Srtgears root folder
const (
	relativeRlsBldFldr = "cmd/srtgears"               // Relative release build folder
	relativeWebFldr    = "web"                        // Relative web folder
	relativeTargetFldr = relativeWebFldr + "/release" // Relative release target folder
)

var (
	srtgearsRoot = "../../" // Srtgears root folder
	targetFolder string     // Target folder to put packed releases into
	rlsBldFldr   string     // Release build folder
)

var targetFiles []string // Target packed releases (*.zip files)

// Mapping from GOOS value to OS display name
var osNameMap = map[string]string{
	"windows": "Windows",
	"linux":   "Linux",
	"darwin":  "OS X",
}

// Mapping from GOARCH value to architecture display name
var archNameMap = map[string]string{
	"amd64": "64-bit",
	"386":   "32-bit",
}

func main() {
	var err error

	if err = initFolders(); err != nil {
		log.Println(err)
		return
	}

	if err = scanAndPack(); err != nil {
		log.Println(err)
		return
	}

	if err = generateDownloadHTML(); err != nil {
		log.Println(err)
		return
	}
}

// initFolders initializes folder variables and checks whether required folders exist.
func initFolders() (err error) {
	if len(os.Args) > 1 {
		srtgearsRoot = os.Args[1]
	}
	if srtgearsRoot, err = filepath.Abs(srtgearsRoot); err != nil {
		return
	}

	rlsBldFldr = filepath.Join(srtgearsRoot, relativeRlsBldFldr)
	targetFolder = filepath.Join(srtgearsRoot, relativeTargetFldr)

	// Check folders:
	for _, folder := range []string{rlsBldFldr, targetFolder} {
		if fi, err := os.Stat(folder); os.IsNotExist(err) {
			return fmt.Errorf("Folder does not exist: %s", folder)
		} else {
			if !fi.IsDir() {
				return fmt.Errorf("Path is not a folder: %s", folder)
			}
		}
	}
	return
}

// scanAndPack scans the release build folder and packs releases.
func scanAndPack() (err error) {
	log.Println("Scanning release build folder:", rlsBldFldr)

	var fis []os.FileInfo
	if fis, err = ioutil.ReadDir(rlsBldFldr); err != nil {
		panic(err)
	}

	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}
		if !rlsFldrPttrn.MatchString(fi.Name()) {
			continue
		}

		var targetFileName string
		if targetFileName, err = packRelease(filepath.Join(rlsBldFldr, fi.Name())); err != nil {
			return
		}
		targetFiles = append(targetFiles, targetFileName)
	}
	return
}

// packRelease packs the content of the specified folder into a zip file, placed under the target folder.
func packRelease(folder string) (targetFileName string, err error) {
	log.Println("Packing release folder:", folder)

	targetFileName = filepath.Join(targetFolder, filepath.Base(folder)+".zip")
	log.Println("To:", targetFileName)

	root := filepath.Dir(folder)

	f, err := os.Create(targetFileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fh, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		// Zip expects forward slashes ('/'), so replace os dependant separators
		relPath = filepath.ToSlash(relPath)
		// fh.Name is only the fienme, so:
		fh.Name = relPath

		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}
		if _, err = w.Write(content); err != nil {
			return err
		}
		return nil
	})

	return
}

// generateDownloadHTML generates the HTML table in the format ready for the Downloads page (download.html).
func generateDownloadHTML() (err error) {
	// It happens we want releases in reverse order in the HTML table:
	sort.Sort(sort.Reverse(sort.StringSlice(targetFiles)))

	webFolder := filepath.Join(srtgearsRoot, relativeWebFldr)

	// hash and print HTML table
	for _, targetFile := range targetFiles {
		url, err := filepath.Rel(webFolder, targetFile)
		if err != nil {
			return err
		}

		// TODO
		// We need forward slashes "/" in urls:
		url = "/" + filepath.ToSlash(url)
		log.Println(url)
	}

	return
}
