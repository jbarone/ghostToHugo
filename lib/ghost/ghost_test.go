package ghost

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"
)

func initLocation() error {
	location, err := time.LoadLocation("UTC")
	if err != nil {
		return err
	}
	SetLocation(location)
	return nil
}

func TestPostIsPage(t *testing.T) {
	for _, d := range boolTestData {
		p := Post{
			Page: json.RawMessage(d.Value),
		}

		if p.IsPage() != d.Want {
			t.Errorf("Expected: %v Actual: %v", d.Want, p.IsPage())
		}
	}
}

func TestPostPublished(t *testing.T) {
	if err := initLocation(); err != nil {
		t.Fatalf("%v", err)
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

	var p Post
	var expected = time.Date(2010, 9, 6, 13, 44, 9, 0, time.UTC)
	for _, d := range testdata {
		SetDateFormat(d.format)
		p.PublishedAt = json.RawMessage(d.value)
		if p.Published() != expected {
			t.Errorf("Parsing %q Expected: %v Actual: %v", d.value, expected, p.Published())
		}
	}
}

func TestPostAuthor(t *testing.T) {
	user := User{
		ID:   1,
		Name: "Test",
	}
	post := Post{
		AuthorID: 1,
	}
	data := &ExportData{
		Users: []User{user},
	}

	if user != *post.Author(data) {
		t.Errorf("Expected: %v Actual: %v", user, post.Author(data))
	}
}

func TestPostAuthorNotFound(t *testing.T) {
	user := User{
		ID:   2,
		Name: "Test",
	}
	post := Post{
		AuthorID: 1,
	}
	data := &ExportData{
		Users: []User{user},
	}

	if nil != post.Author(data) {
		t.Errorf("Expected: %v Actual: %v", nil, post.Author(data))
	}
}

func TestPostAuthorNullData(t *testing.T) {
	post := Post{
		AuthorID: 1,
	}

	if nil != post.Author(nil) {
		t.Errorf("Expected: %v Actual: %v", nil, post.Author(nil))
	}
}

func TestPostTags(t *testing.T) {
	tag := Tag{ID: 1}
	post := Post{
		ID: 1,
	}
	data := &ExportData{
		Tags:     []Tag{tag},
		PostTags: []PostTag{{PostID: 1, TagID: 1}},
	}

	tags := post.Tags(data)
	if len(tags) != 1 {
		t.Errorf("Expected: 1 tag  Actual: %v tag(s)", len(tags))
	} else if tags[0].ID != tag.ID {
		t.Errorf("Expected: %v Actual: %v", tag.ID, tags[0].ID)
	}
}

func TestProcessWrapped(t *testing.T) {
	data, err := os.Open("testdata/wrapped.json")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer data.Close()

	reader := ExportReader{
		ReadSeeker: data,
	}

	entries, err := processWrapped(reader)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected: 1 entry  Actual: %v entries", len(entries))
	} else if reflect.DeepEqual(exportRecord, entries[0]) {
		t.Errorf("Expected: %v Actual: %v", exportRecord, entries[0])
	}
}

func TestProcessUnwrapped(t *testing.T) {
	data, err := os.Open("testdata/unwrapped.json")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer data.Close()

	reader := ExportReader{
		ReadSeeker: data,
	}

	entries, err := processUnwrapped(reader)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected: 1 entry  Actual: %v entries", len(entries))
	} else if reflect.DeepEqual(exportRecord, entries[0]) {
		t.Errorf("Expected: %v Actual: %v", exportRecord, entries[0])
	}
}

func TestProcess(t *testing.T) {
	for _, filename := range []string{
		"testdata/wrapped.json",
		"testdata/unwrapped.json",
	} {
		data, err := os.Open(filename)
		if err != nil {
			t.Fatalf("%v", err)
		}
		defer data.Close()
		reader := ExportReader{
			ReadSeeker: data,
		}
		entries, err := Process(reader)
		if err != nil {
			t.Fatalf("%v", err)
		}

		if len(entries) != 1 {
			t.Errorf("Expected: 1 entry  Actual: %v entries", len(entries))
		} else if reflect.DeepEqual(exportRecord, entries[0]) {
			t.Errorf("Expected: %v Actual: %v", exportRecord, entries[0])
		}
	}
}

var boolTestData = []struct {
	Value string
	Want  bool
}{
	{"1", true},
	{"0", false},
	{"true", true},
	{"false", false},
	{"nonsense", false},
}

var exportRecord = ExportEntry{
	Meta: ExportMeta{
		ExportedOn: 1388805572000,
		Version:    "003",
	},
	Data: ExportData{
		Posts: []Post{
			{
				ID:          5,
				Title:       "my blog post title",
				Slug:        "my-blog-post-title",
				Content:     "the *markdown* formatted post body",
				Status:      "published",
				AuthorID:    1,
				PublishedAt: json.RawMessage("1283780649000"),
			},
		},
		Tags: []Tag{
			{
				ID:   3,
				Name: "Colorado Ho!",
			},
			{
				ID:   4,
				Name: "blue",
			},
		},
		PostTags: []PostTag{
			{
				TagID:  3,
				PostID: 5,
			},
			{
				TagID:  3,
				PostID: 2,
			},
			{
				TagID:  4,
				PostID: 24,
			},
		},
		Users: []User{
			{
				ID:   2,
				Name: "user's name",
			},
		},
	},
}
