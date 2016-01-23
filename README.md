# Srtgears&trade;

Srtgears&trade; is a subtitle engine for reading subtitle files, manipulating / transforming them and then saving the result into another file.

Srtgears provides some very handy features which are not available in other subtitle tools, for example:

- merge 2 subtitle files to have dual subs: one at the bottom, one at the top (this is not concatenation, but that's also supported)
- lengthen / shorten display duration of subtitles (if you're a slow reader, you're gonna appreciate this :))
- remove hearing impact-only subtitles (such as "[PHONE RINGING]" or "(phone ringing)")
- strip off formatting (such as &lt;i&gt;, &lt;b&gt;, &lt;u&gt;, &lt;font&gt;) 
- statistics from the subtitle
- etc...

Home page: https://srt-gears.appspot.com

# Presentation

The Srtgears engine is presented in 3 ways:

## 1. Command line tool
Srtgears is available as a command line tool for easy and fast usage.

Binary (compiled) distributions are available on the download page:

https://srt-gears.appspot.com/downloads

## 2. Web interface: online web page

Srtgears can also be used on the web for those who do not want to download the tool, just try it or use it from the browser. It can be found here:

https://srt-gears.appspot.com/srtgears-online

## 3. Go package

And last (but not least) a Go package with clear and easy-to use interface + documentation, should you need to use it in your project.

To get the source code (along with the sources of the tool and web interface), use `go get`:

    go get github.com/gophergala2016/srtgears

# LICENSE

See [LICENSE](https://github.com/gophergala2016/srtgears/blob/master/LICENSE.md)
