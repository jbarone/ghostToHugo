// Copyright © 2016-present Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gitmap

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

var (
	revision   = "4d9ad733fa40310607ebe9f700d59dcac93ace89"
	repository string
)

func init() {
	var err error
	if repository, err = os.Getwd(); err != nil {
		panic(err)
	}
}

func TestMap(t *testing.T) {
	var (
		gm  GitMap
		gr  *GitRepo
		err error
	)

	if gr, err = Map(repository, revision); err != nil {
		t.Fatal(err)
	}

	gm = gr.Files

	if len(gm) != 9 {
		t.Fatalf("Wrong number of files, got %d, expected %d", len(gm), 9)
	}

	assertFile(t, gm,
		"testfiles/d1/d1.txt",
		"4d9ad73",
		"4d9ad733fa40310607ebe9f700d59dcac93ace89",
		"2016-07-20",
	)

	assertFile(t, gm,
		"testfiles/d2/d2.txt",
		"9d1dc47",
		"9d1dc478eef267829831226d913a3ca249c489d4",
		"2016-07-19",
	)

	assertFile(t, gm,
		"README.md",
		"866cbcc",
		"866cbccdab588b9908887ffd3b4f2667e94090c3",
		"2016-07-20",
	)
}

func assertFile(
	t *testing.T,
	gm GitMap,
	filename,
	expectedAbbreviatedHash,
	expectedHash,
	expectedDate string) {

	var (
		gi *GitInfo
		ok bool
	)

	if gi, ok = gm[filename]; !ok {
		t.Fatal(filename)
	}

	if gi.AbbreviatedHash != expectedAbbreviatedHash || gi.Hash != expectedHash {
		t.Error("Invalid tree hash, file", filename, "abbreviated:", gi.AbbreviatedHash, "full:", gi.Hash, gi.Subject)
	}

	if gi.AuthorName != "Bjørn Erik Pedersen" || gi.AuthorEmail != "bjorn.erik.pedersen@gmail.com" {
		t.Error("These commits are mine! Got", gi.AuthorName, "and", gi.AuthorEmail)
	}

	if gi.AuthorDate.Format("2006-01-02") != expectedDate {
		t.Error("Invalid date:", gi.AuthorDate)
	}
}

func TestActiveRevision(t *testing.T) {
	var (
		gm  GitMap
		gr  *GitRepo
		err error
	)

	if gr, err = Map(repository, "HEAD"); err != nil {
		t.Fatal(err)
	}

	gm = gr.Files

	if len(gm) < 10 {
		t.Fatalf("Wrong number of files, got %d, expected at least %d", len(gm), 10)
	}

	if len(gm) < 10 {
		t.Fatalf("Wrong number of files, got %d, expected at least %d", len(gm), 10)
	}
}

func TestGitExecutableNotFound(t *testing.T) {
	defer initDefaults()
	gitExec = "thisShouldHopefullyNotExistOnPath"
	gi, err := Map(repository, revision)

	if err != GitNotFound || gi != nil {
		t.Fatal("Invalid error handling")
	}

}

func TestEncodeJSON(t *testing.T) {
	var (
		gm       GitMap
		gr       *GitRepo
		gi       *GitInfo
		err      error
		ok       bool
		filename = "README.md"
	)

	if gr, err = Map(repository, revision); err != nil {
		t.Fatal(err)
	}

	gm = gr.Files

	if gi, ok = gm[filename]; !ok {
		t.Fatal(filename)
	}

	b, err := json.Marshal(&gi)

	if err != nil {
		t.Fatal(err)
	}

	s := string(b)

	if s != `{"hash":"866cbccdab588b9908887ffd3b4f2667e94090c3","abbreviatedHash":"866cbcc","subject":"Add codecov to Travis config","authorName":"Bjørn Erik Pedersen","authorEmail":"bjorn.erik.pedersen@gmail.com","authorDate":"2016-07-20T02:22:23+02:00"}` {
		t.Errorf("JSON marshal error: \n%s", s)
	}
}

func TestGitRevisionNotFound(t *testing.T) {
	gi, err := Map(repository, "adfasdfasdf")

	// TODO(bep) improve error handling.
	if err == nil || gi != nil {
		t.Fatal("Invalid error handling", err)
	}
}

func TestGitRepoNotFound(t *testing.T) {
	gi, err := Map("adfasdfasdf", revision)

	// TODO(bep) improve error handling.
	if err == nil || gi != nil {
		t.Fatal("Invalid error handling", err)
	}
}

func TestTopLevelAbsPath(t *testing.T) {
	var (
		gr  *GitRepo
		err error
	)

	if gr, err = Map(repository, revision); err != nil {
		t.Fatal(err)
	}

	expected := "github.com/bep/gitmap"

	if !strings.HasSuffix(gr.TopLevelAbsPath, expected) {
		t.Fatalf("Expected to end with %q got %q", expected, gr.TopLevelAbsPath)
	}
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Map(repository, revision)
		if err != nil {
			b.Fatalf("Got error: %s", err)
		}
	}
}
