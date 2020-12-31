# ghostToHugo

[![GitHub release](https://img.shields.io/github/release/jbarone/ghostToHugo.svg)](https://github.com/jbarone/ghostToHugo/releases/latest)
[![MIT License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](/LICENSE)
[![Build Status](https://travis-ci.org/jbarone/ghostToHugo.svg?branch=master)](https://travis-ci.org/jbarone/ghostToHugo)
[![Go Report Card](https://goreportcard.com/badge/github.com/jbarone/ghostToHugo)](https://goreportcard.com/report/github.com/jbarone/ghostToHugo)

**ghostToHugo** is a utility project that was created the allow the conversion
of an export from the Ghost blogging engine into the Hugo engine.

## Installing

There are 2 options for installation.

### Pre-Built Binaries

With every new versioned release. Binaries are built for most major platforms.
You can simply download the one for your operating system from the
[releases page](https://github.com/jbarone/ghostToHugo/releases/latest). Unzip
the package and place it somewhere in your path, and you are ready to go.

### Fetch With Go

The project is written in Go, and can easily be built from source.
Make sure you have Go installed and configured, then just run:

```
GO111MODULE=on go get github.com/jbarone/ghostToHugo@latest
```

this will download, compile, and install `ghostToHugo`

### Building the Latest and Greatest

Alternately, you can build directly:

```
git clone http://github.com/jbarone/ghostToHugo
cd ghostToHugo
go install
```

## Usage

```
Usage: ghostToHugo [OPTIONS] <Ghost Export>
  -d, --dateformat string   date format string to use for time conversions (default "2006-01-02 15:04:05")
      --debug               print verbose logging output
  -f, --force               allow import into non-empty target directory
  -p, --hugo string         path to create the new hugo project (default "newhugosite")
  -l, --location string     location to use for time conversions (default: local)
  -v, --verbose             print verbose logging output
```

At a minimum you need to specify the path to the exported Ghost json file.

NOTES:

- The `dateformat` string must be provided in Go's specific time format string. Reference [here](https://gobyexample.com/time-formatting-parsing)
- The `location` string should be a value that matches the IANA Time Zone database, such as "America/New_York"
- The path specified for the new Hugo site, must either not exist, or be an empty directory. A new site will be created at that location.

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
$ ghostToHugo -d "2006-01-02 15:04:05" export.json
```

```
$ ghostToHugo --location "America/Chicago" export.json
$ ghostToHugo -l "America/Chicago" export.json
```

## Exporting your Ghost content
You can export your Ghost content (and settings) from the "Labs" section of your Ghost install, which will be at a URL like:

```http(s)://<your blog url>/ghost/settings/labs/```

See this [Ghost support article](https://help.ghost.org/hc/en-us/articles/224112927-Import-Export-Data) for details.
