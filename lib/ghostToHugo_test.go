package ghostToHugo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	location, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatal(err)
	}

	var testdata = []struct {
		format string
		value  string
	}{
		{time.RFC3339, "1283780649000"},
		{time.RFC3339, `"2010-09-06T13:44:09-00:00"`},
		{"2006-01-02T15:04:05", `"2010-09-06T13:44:09"`},
		{"2006-01-02 15:04:05", `"2010-09-06 13:44:09"`},
	}

	var expected = time.Date(2010, 9, 6, 13, 44, 9, 0, time.UTC)
	for _, data := range testdata {
		gth, err := NewGhostToHugo(
			WithLocation(location),
			WithDateFormat(data.format),
		)
		if err != nil {
			t.Fatal(err)
		}

		result := gth.parseTime(json.RawMessage(data.value))

		if result != expected {
			t.Errorf("Parsing %q Expected: %v Actual: %v",
				data.value, expected, result)
		}
	}
}

func TestImportGhost(t *testing.T) {
	data := []string{
		"wrapped.json",
		"unwrapped.json",
	}
	for _, d := range data {
		f, err := os.Open(filepath.Join("testdata", d))
		if err != nil {
			t.Fatal(err)
		}

		gth, _ := NewGhostToHugo()
		posts, err := gth.importGhost(f)
		if err != nil {
			t.Fatal(err)
		}

		var entryCount int
		for p := range posts {
			entryCount++
			if reflect.DeepEqual(expectedPost, p) {
				t.Errorf("Expected: %v Actual: %v", expectedPost, p)
			}
		}

		if entryCount != 1 {
			t.Errorf("Expected 1 entry, found %v", entryCount)
		}

		err = f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
}

var expectedPost = Post{
	ID:              5,
	Title:           "my blog post title",
	Slug:            "my-blog-post-title",
	Content:         "the *markdown* formatted post body",
	Image:           "",
	Page:            json.RawMessage("0"),
	Status:          "published",
	MetaDescription: "",
	AuthorID:        1,
	PublishedAt:     json.RawMessage("1283780649000"),
	CreatedAt:       json.RawMessage("1283780649000"),

	IsDraft: false,
	IsPage:  false,
	Author:  "user's name",
	Tags:    []string{"Colorado Ho!"},
}
