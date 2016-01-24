# Srtgears&trade;

Srtgears&trade; is a subtitle engine for reading subtitle files, manipulating / transforming them and then saving the result into another file.

Srtgears provides some very handy features which are not available in other subtitle tools, for example:

- merge 2 subtitle files to have dual subs: one at the bottom, one at the top (this is not concatenation, but that's also supported)
- lengthen / shorten display duration of subtitles (if you're a slow reader, you're gonna appreciate this :))
- remove hearing impaired texts (such as "[PHONE RINGING]" or "(phone ringing)")
- strip off formatting (such as &lt;i&gt;, &lt;b&gt;, &lt;u&gt;, &lt;font&gt;) 
- statistics from the subtitle
- etc...

Home page: https://srt-gears.appspot.com

## Presentation

The Srtgears engine is presented in 3 ways:

### 1. Command line tool
Srtgears is available as a command line tool for easy, fast, scriptable and repeatable usage.

Binary (compiled) distributions are available on the download page:

https://srt-gears.appspot.com/download.html

The command line tool uses only the Go standard library and the srtgears engine (see below).

### 2. Web interface: online web page

Srtgears can also be used on the web for those who do not want to download the tool just it out from the browser. It can be found here:

https://srt-gears.appspot.com/srtgears-online.html

The web interface is a Google App Engine project, implemented using the Go AppEngine SDK. The server side of the web interface uses the srtgears engine (see below).

The `web` folder is the root of the App Engine project. See instructions below if you want to fire up the web interface locally.

### 3. Srtgears engine: a Go package

And last (but not least) a Go package for developers. The engine was designed to be independent from the command line and web interfaces, its API is clear, well documented and easy-to-use.

To get the source code (along with the sources of the tool and web interface), use `go get`:

    go get github.com/gophergala2016/srtgears
    
Documentation can be found at:

http://godoc.org/github.com/gophergala2016/srtgears

To use the engine, first import it:

    import "github.com/gophergala2016/srtgears"

And for example using the engine to merge 2 subtitle files to have a dual sub saved in Sub Station Alpha (*.ssa) format:

	sp1, err := srtgears.ReadSrtFile("eng.srt")
	// check err
	sp2, err := srtgears.ReadSrtFile("hun.srt")
	// check err
	sp1.Merge(sp2)
	err = srtgears.WriteSsaFile("eng+hun.ssa", sp1);
	// check err

You can see more usage examples in the [package doc](http://godoc.org/github.com/gophergala2016/srtgears).

Also worth noting that the subtitle transformations of the command line tool and the web interface are driven by the same `Executor`, it is "outsourced" to the `github.com/gophergala2016/srtgears/exec` package.

**Instructions to fire up the web interface locally:**

If you want to try it locally, you need to download the Go AppEngine SDK. Once it's done, the srtgears engine (`srtgears` package) and the `Executor` (`srtgears/exec` package) must be copied into the web folder (it is not committed here to avoid duplication). On Windows this can be done by running the `copy-srtgears-here.cmd` script in the web folder.

If this is ready, it can be started locally by running the `goapp serve` command of the SDK from the `web` folder.

## Limits

Input files must be UTF-8 encoded, output files will be UTF-8 encoded as well.

Supported input format is SubRip (`*.srt`) only, supported output formats are SubRip (`*.srt`) and Sub Station Alpha (`*.ssa`).

It should also be noted that SubRip format specification does not include subtitle positioning. Srtgears uses an unofficial extension `{\anX}` which may not be supported by all video players, or some players interpret the position values differently. [MPC-HC](https://mpc-hc.org/) has full support for it. In these cases the Sub Station Alpha output format is recommended (where the specification covers subtitle positioning / alignment).

## License

See [LICENSE](https://github.com/gophergala2016/srtgears/blob/master/LICENSE.md)
