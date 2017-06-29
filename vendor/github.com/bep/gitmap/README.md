# GitMap

[![GoDoc](https://godoc.org/github.com/bep/gitmap?status.svg)](https://godoc.org/github.com/bep/gitmap)
[![Build Status](https://travis-ci.org/bep/gitmap.svg)](https://travis-ci.org/bep/gitmap) [![Build status](https://ci.appveyor.com/api/projects/status/c8tu1wdoa4j7q81g?svg=true)](https://ci.appveyor.com/project/bjornerik/gitmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/bep/gitmap)](https://goreportcard.com/report/github.com/bep/gitmap)
[![codecov](https://codecov.io/gh/bep/gitmap/branch/master/graph/badge.svg)](https://codecov.io/gh/bep/gitmap)

A fairly fast way to create a map from all the filenames to info objects for a given revision of a Git repo.

This library uses `os/exec` to talk to Git. There are faster ways to do this by using some Go Git-lib or C bindings, but that adds dependencies I really don't want or need.

If some `git log kung fu master` out there have suggestions for improvements, please open an issue or a PR.
