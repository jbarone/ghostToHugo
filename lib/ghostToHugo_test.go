package ghostToHugo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestImportGhost(t *testing.T) {
	data := []string{
		"wrapped.json",
		"unwrapped.json",
	}
	for _, d := range data {
		f, err := os.Open(filepath.Join("testdata", d))
		if err != nil {
			t.Error(err)
		}

		gth, _ := NewGhostToHugo()
		posts, err := gth.importGhost(f)
		if err != nil {
			t.Error(err)
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
			t.Error(err)
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
