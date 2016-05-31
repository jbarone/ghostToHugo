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

func TestPostIsFeatured(t *testing.T) {
	for _, d := range boolTestData {
		p := Post{
			Featured: json.RawMessage(d.Value),
		}

		if p.IsFeatured() != d.Want {
			t.Errorf("Expected: %v Actual: %v", d.Want, p.IsFeatured())
		}
	}
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

func TestPostCreated(t *testing.T) {
	if err := initLocation(); err != nil {
		t.Fatalf("%v", err)
	}
	p := Post{
		CreatedAt: 1283780649000,
	}

	var expected = time.Date(2010, 9, 6, 13, 44, 9, 0, time.UTC)
	if p.Created() != expected {
		t.Errorf("Expected: %v Actual: %v", expected, p.Created())
	}
}

func TestPostUpdated(t *testing.T) {
	if err := initLocation(); err != nil {
		t.Fatalf("%v", err)
	}
	p := Post{
		UpdatedAt: 1286958624000,
	}

	var expected = time.Date(2010, 10, 13, 8, 30, 24, 0, time.UTC)
	if p.Updated() != expected {
		t.Errorf("Expected: %v Actual: %v", expected, p.Updated())
	}
}

func TestPostPublished(t *testing.T) {
	if err := initLocation(); err != nil {
		t.Fatalf("%v", err)
	}
	p := Post{
		PublishedAt: 1283780649000,
	}

	var expected = time.Date(2010, 9, 6, 13, 44, 9, 0, time.UTC)
	if p.Published() != expected {
		t.Errorf("Expected: %v Actual: %v", expected, p.Published())
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
		PostTags: []PostTag{PostTag{PostID: 1, TagID: 1}},
	}

	tags := post.Tags(data)
	if len(tags) != 1 {
		t.Errorf("Expected: 1 tag  Actual: %v tag(s)", len(tags))
	} else if tags[0].ID != tag.ID {
		t.Errorf("Expected: %v Actual: %v", tag.ID, tags[0].ID)
	}
}

func TestUserCreated(t *testing.T) {
	if err := initLocation(); err != nil {
		t.Fatalf("%v", err)
	}
	u := User{
		CreatedAt: 1283780649000,
	}

	var expected = time.Date(2010, 9, 6, 13, 44, 9, 0, time.UTC)
	if u.Created() != expected {
		t.Errorf("Expected: %v Actual: %v", expected, u.Created())
	}
}

func TestUserUpdated(t *testing.T) {
	if err := initLocation(); err != nil {
		t.Fatalf("%v", err)
	}
	u := User{
		UpdatedAt: 1286958624000,
	}

	var expected = time.Date(2010, 10, 13, 8, 30, 24, 0, time.UTC)
	if u.Updated() != expected {
		t.Errorf("Expected: %v Actual: %v", expected, u.Updated())
	}
}

func TestTagIsHidden(t *testing.T) {
	for _, d := range boolTestData {
		tag := Tag{
			Hidden: json.RawMessage(d.Value),
		}

		if tag.IsHidden() != d.Want {
			t.Errorf("Expected: %v Actual: %v", d.Want, tag.IsHidden())
		}
	}
}

func TestTagCreated(t *testing.T) {
	if err := initLocation(); err != nil {
		t.Fatalf("%v", err)
	}
	p := Tag{
		CreatedAt: 1283780649000,
	}

	var expected = time.Date(2010, 9, 6, 13, 44, 9, 0, time.UTC)
	if p.Created() != expected {
		t.Errorf("Expected: %v Actual: %v", expected, p.Created())
	}
}

func TestTagUpdated(t *testing.T) {
	if err := initLocation(); err != nil {
		t.Fatalf("%v", err)
	}
	p := Tag{
		UpdatedAt: 1286958624000,
	}

	var expected = time.Date(2010, 10, 13, 8, 30, 24, 0, time.UTC)
	if p.Updated() != expected {
		t.Errorf("Expected: %v Actual: %v", expected, p.Updated())
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
			Post{
				ID:          5,
				Title:       "my blog post title",
				Slug:        "my-blog-post-title",
				Content:     "the *markdown* formatted post body",
				HTML:        "the <i>html</i> formatted post body",
				Status:      "published",
				Language:    "en_US",
				AuthorID:    1,
				CreatedAt:   1283780649000,
				CreatedBy:   1,
				UpdatedAt:   1286958624000,
				UpdatedBy:   1,
				PublishedAt: 1283780649000,
				PublishedBy: 1,
			},
		},
		Tags: []Tag{
			Tag{
				ID:   3,
				Name: "Colorado Ho!",
				Slug: "colorado-ho",
			},
			Tag{
				ID:   4,
				Name: "blue",
				Slug: "blue",
			},
		},
		PostTags: []PostTag{
			PostTag{
				TagID:  3,
				PostID: 5,
			},
			PostTag{
				TagID:  3,
				PostID: 2,
			},
			PostTag{
				TagID:  4,
				PostID: 24,
			},
		},
		Users: []User{
			User{
				ID:        2,
				Name:      "user's name",
				Slug:      "users-name",
				Email:     "user@example.com",
				Status:    "active",
				Language:  "en_US",
				CreatedAt: 1283780649000,
				CreatedBy: 1,
				UpdatedAt: 1286958624000,
				UpdatedBy: 1,
			},
		},
	},
}
