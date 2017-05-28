package ghostToHugo

import (
	"reflect"
	"testing"
	"time"
)

// func TestPath(t *testing.T) {
// 	data := []struct {
// 		post post
// 		path string
// 	}{
// 		{
// 			post{
// 				Slug:   "test",
// 				IsPage: true,
// 			},
// 			path.Join("content", "test.md"),
// 		},
// 		{
// 			post{
// 				Slug:   "test",
// 				IsPage: false,
// 			},
// 			path.Join("content", "post", "test.md"),
// 		},
// 	}

// 	for _, d := range data {
// 		if d.path != d.post.path() {
// 			t.Errorf("Post.path() = %s; want %s", d.post.path(), d.path)
// 		}
// 	}
// }

func TestMetadata(t *testing.T) {
	data := []struct {
		post     post
		expected map[string]interface{}
	}{
		{
			post{
				Title:           "Title",
				Slug:            "test",
				IsDraft:         true,
				MetaDescription: "test description",
				Created:         time.Now(),
			},
			map[string]interface{}{
				"title":       "Title",
				"draft":       true,
				"slug":        "test",
				"description": "test description",
			},
		},
		{
			post{
				Title:           "Title",
				Slug:            "test",
				IsDraft:         false,
				MetaDescription: "test description",
				Published:       time.Now(),
			},
			map[string]interface{}{
				"title":       "Title",
				"draft":       false,
				"slug":        "test",
				"description": "test description",
			},
		},
		{
			post{
				Title:           "Title",
				Slug:            "test",
				IsDraft:         false,
				MetaDescription: "test description",
				Published:       time.Now(),
				Image:           "/content/test.jpg",
			},
			map[string]interface{}{
				"title":       "Title",
				"draft":       false,
				"slug":        "test",
				"description": "test description",
				"image":       "/test.jpg",
			},
		},
		{
			post{
				Title:           "Title",
				Slug:            "test",
				IsDraft:         false,
				MetaDescription: "test description",
				Published:       time.Now(),
				Tags:            []string{"Test"},
			},
			map[string]interface{}{
				"title":       "Title",
				"draft":       false,
				"slug":        "test",
				"description": "test description",
				"tags":        []string{"Test"},
				"categories":  []string{"Test"},
			},
		},
		{
			post{
				Title:           "Title",
				Slug:            "test",
				IsDraft:         false,
				MetaDescription: "test description",
				Published:       time.Now(),
				Author:          "Aurthor",
			},
			map[string]interface{}{
				"title":       "Title",
				"draft":       false,
				"slug":        "test",
				"description": "test description",
				"author":      "Aurthor",
			},
		},
	}

	for _, d := range data {
		if reflect.DeepEqual(d.expected, d.post.metadata()) {
			t.Errorf(
				"Post.metadata() = %v; want %v",
				d.post.metadata(),
				d.expected,
			)
		}
	}
}
