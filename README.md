# ghostToHugo

[![GitHub release](https://img.shields.io/github/release/tony19/x-user-dropdown.svg)](https://github.com/jbarone/ghostToHugo/releases/latest)
[![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE)
[![Build Status](https://travis-ci.org/jbarone/ghostToHugo.svg?branch=master)](https://travis-ci.org/jbarone/ghostToHugo)
[![Go Report Card](https://goreportcard.com/badge/github.com/jbarone/ghostToHugo)](https://goreportcard.com/report/github.com/jbarone/ghostToHugo)

**ghostToHugo** is a utility project that was created the allow the conversion
of an export from the Ghost blogging engine into the Hugo engine.

## Installing

The project is written in Go, and currently will require building from source.
Make sure you have Go installed and configured, then just run:

```
go get -u github.com/jbarone/ghostToHugo
```

this will download, compile, and install `ghostToHugo`

## Using

```
Usage: ghostToHugo [OPTIONS] <Ghost Export>
  -f, --dateformat string   Date format string to use for time conversions (default: RFC3339)
  -p, --hugo string         Path to hugo project (default ".")
  -l, --location string     Location to use for time conversions (default: local)
```

At a minimum you need to specify the path to the exported Ghost json file.

NOTES: 

- The `dateformat` string should be provided in Go's time format string. Reference [here](https://golang.org/src/time/format.go)
- The `location` string should be a value that matches the IANA Time Zone database, such as "America/New_York"

### Examples

```
$ ghostToHugo export.json
```

```
$ ghostToHugo --hugo ~/mysite export.json
$ ghostToHugo -p ~/mysite export.json
```

```
$ ghostToHugo --dateformat "2006-01-02 15:04:05" export.json
$ ghostToHugo -f "2006-01-02 15:04:05" export.json
```

```
$ ghostToHugo --location "America/Chicago" export.json
$ ghostToHugo -l "America/Chicago" export.json
```

## Exporting your Ghost content
You can export your Ghost content (and settings) from the "Labs" section of your Ghost install, which will be at a URL like `<your blog url>/ghost/settings/labs/`.

See this [Ghost support article](https://help.ghost.org/hc/en-us/articles/224112927-Import-Export-Data) for details.
