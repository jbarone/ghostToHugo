package hugo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/jbarone/ghostToHugo/lib/ghost"
	"github.com/spf13/viper"
)

func TestGetPath(t *testing.T) {
	initViper()

	var testdata = []struct {
		post ghost.Post
		want string
	}{
		{
			post: ghost.Post{Slug: "test", Page: json.RawMessage([]byte{'1'})},
			want: filepath.Join(viper.GetString("contentDir"), "test.md"),
		},
		{
			post: ghost.Post{Slug: "test", Page: json.RawMessage([]byte{'0'})},
			want: filepath.Join(viper.GetString("contentDir"), "post", "test.md"),
		},
	}

	for _, data := range testdata {
		s := getPath(data.post)
		if s != data.want {
			t.Errorf("Expected: %v  Actual: %v", data.want, s)
		}
	}
}

func TestTagsToStringSlice(t *testing.T) {
	testdata := []struct {
		tags []ghost.Tag
		want []string
	}{
		{
			tags: nil,
			want: nil,
		},
		{
			tags: []ghost.Tag{},
			want: nil,
		},
		{
			tags: []ghost.Tag{
				ghost.Tag{Name: "test"},
			},
			want: []string{"test"},
		},
		{
			tags: []ghost.Tag{
				ghost.Tag{Name: "test1"},
				ghost.Tag{Name: "test2"},
			},
			want: []string{
				"test1",
				"test2",
			},
		},
	}

	for _, data := range testdata {
		s := tagsToSlice(data.tags)
		if len(s) != len(data.want) {
			t.Errorf("Length Mismatch Expected: %v  Actual: %v", len(data.want), len(s))
			continue
		}
		for i := range data.want {
			if s[i] != data.want[i] {
				t.Errorf("Expected: %v  Actual: %v", data.want[i], s[i])
				break
			}
		}
	}
}

func initViper() {
	viper.Reset()
	viper.Set("MetaDataFormat", "toml")
	viper.Set("contentDir", filepath.Join(os.TempDir(), "content"))
}

// func initFs() error {
// 	hugofs.SetSource(new(afero.MemMapFs))
// 	perm := os.FileMode(0755)
// 	var err error

// 	// create directories
// 	dir = filepath.Join(os.TempDir(), "content")
// 	err = hugofs.Source().Mkdir(dir, perm)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
