# Srtgears&trade;

Srtgears&trade; is a subtitle engine for reading subtitle files, manipulating / transforming them and then saving the result into another file.

Srtgears provides some very handy features which are not available in other subtitle tools, for example:

- merge 2 subtitle files to have dual subs (one at the bottom, one at the top)
- lengthen / shorten display duration of subtitles (if you're a slow reader, you're gonna appreciate this :))
- remove hearing impact-only subtitles (such as "[PHONE RINGING]" or "(phone ringing)")
- strip off formatting (such as &lt;i&gt;, &lt;b&gt;, &lt;u&gt;, &lt;font&gt;) 
- statistics from the subtitle
- etc...

Home page: https://srt-gears.appspot.com

# Presentation

The Srtgears engine is presented in 3 ways:

- a command line tool for easy and fast usage
- a web page for those who do not want to download the tool, just try it or use it from the browser
- a Go package with clear and easy-to use interface + documentation, should you need to use it in your project

# Srtgears online

The engine is available online as a webpage here:

https://srt-gears.appspot.com/srtgears-online


# How to get it or install it

Binary (compiled) distributions are available on the download page:

https://srt-gears.appspot.com/downloads

To get the source code, use `go get`:

    go get github.com/gophergala2016/srtgears

# LICENSE

See [LICENSE](https://github.com/gophergala2016/srtgears/blob/master/LICENSE.md)
