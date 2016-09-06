/*

This is the main package of the packrelease tool which packs cross-compiled releases
into the web/release folder, and also generates the HTML table in the format ready for the
download.html.

*/
package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

// Release folder name pattern, e.g. "srtgears-1.1-windows-amd64"
var rlsFldrPttrn = regexp.MustCompile("^srtgears-([^-]*)-([^-]*)-([^-]*)$")

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
		fi, err := os.Stat(folder)
		if os.IsNotExist(err) {
			return fmt.Errorf("Folder does not exist: %s", folder)
		}
		if !fi.IsDir() {
			return fmt.Errorf("Path is not a folder: %s", folder)
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

		fh.Method = zip.Deflate
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

// Describes a row (a file) in the download table.
type fileDesc struct {
	Class  string // Css class of the row
	OS     string // OS
	Arch   string // Architecture
	URL    string // Download URL
	Name   string // File name
	Size   string // File size
	SHA256 string // File SHA256 checksum
}

// generateDownloadHTML generates the HTML table in the format ready for the Downloads page (download.html).
func generateDownloadHTML() (err error) {
	fds := []*fileDesc{}

	// It happens we want releases in reverse order in the HTML table:
	sort.Sort(sort.Reverse(sort.StringSlice(targetFiles)))
	webFolder := filepath.Join(srtgearsRoot, relativeWebFldr)

	params := map[string]interface{}{
		"ReleaseDate": time.Now().Format("2006-01-02"),
	}

	// Fill fds scice
	for i, targetFile := range targetFiles {
		fd := fileDesc{Name: filepath.Base(targetFile)}
		fds = append(fds, &fd)

		// the regexp patter is for folder name (without extension)
		nameNoExt := fd.Name[:len(fd.Name)-len(filepath.Ext(fd.Name))]
		if parts := rlsFldrPttrn.FindStringSubmatch(nameNoExt); len(parts) > 0 {
			// [full string, version, os, arch]
			params["Version"] = parts[1]
			fd.OS = osNameMap[parts[2]]
			fd.Arch = archNameMap[parts[3]]
		} else {
			// Never to happen, file name was already matched earlier
			return fmt.Errorf("Target name does not match pattern: %s", targetFile)
		}
		if i%2 != 0 {
			fd.Class = "alt"
		}
		if fd.URL, err = filepath.Rel(webFolder, targetFile); err != nil {
			return
		}
		// We need forward slashes "/" in urls:
		fd.URL = "/" + filepath.ToSlash(fd.URL)
		var fi os.FileInfo
		if fi, err = os.Stat(targetFile); err != nil {
			return
		}
		fd.Size = fmt.Sprintf("%.2f MB", float64(fi.Size())/(1<<20))

		// Hash and include checksum
		var content []byte
		if content, err = ioutil.ReadFile(targetFile); err != nil {
			return
		}
		fd.SHA256 = fmt.Sprintf("%x", sha256.Sum256(content))
	}
	params["Fds"] = fds

	// Now generate download table:
	t := template.Must(template.New("").Parse(dltable))
	buf := &bytes.Buffer{}
	if err = t.Execute(buf, params); err != nil {
		return
	}
	outf := "download-table.html"
	if err = ioutil.WriteFile(outf, buf.Bytes(), 0); err != nil {
		return
	}
	log.Println("Download table written to:", outf)
	// Also print to console:
	os.Stdout.Write(buf.Bytes())
	return
}

const dltable = `			<h4>Latest version: Srtgears {{.Version}}, release date: {{.ReleaseDate}}</h4>
			
			<table class="dlTable">
				<tr>
					<th>OS</th>
					<th>Arch</th>
					<th>Link</th>
					<th>Size</th>
					<th>SHA256 Checksum</th>
				</tr>
{{range $i, $fd := .Fds}}				<tr{{with .Class}} class="{{.}}"{{end}}>
					<td>{{.OS}}</td>
					<td>{{.Arch}}</td>
					<td><a href="{{.URL}}">{{.Name}}</a></td>
					<td>{{.Size}}</td>
					<td class="checksum">{{.SHA256}}</td>
				</tr>
{{end}}			</table>
`
