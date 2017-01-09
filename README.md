# ghostToHugo

**ghostToHugo** is a utility project that was created the allow the conversion
of an export from the Ghost blogging engine into the Hugo engine.

## Installing

The project is written in Go, and currently will require building from source.
Make sure you have Go installed and configured, then just run:

```
go get github.com/jbarone/ghostToHugo
```

this will download, compile, and install `ghostToHugo`

## Using

```
Usage: ghostToHugo [OPTIONS] <Ghost Export>
  -dateformat string
        Date format string to use for time conversions (default: RFC3339)
  -hugo string
        Path to hugo project (default ".")
  -location string
        Location to use for time conversions (default: local)
```

At a minimum you need to specify the path to the exported Ghost json file.

## Exporting your Ghost content
You can export your Ghost content (and settings) from the "Labs" section of your Ghost install, which will be at a URL like `<your blog url>/ghost/settings/labs/`.

See this [Ghost support article](https://help.ghost.org/hc/en-us/articles/224112927-Import-Export-Data) for details.
