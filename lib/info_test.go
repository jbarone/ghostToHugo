package ghostToHugo

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestDecodeGhostInfo(t *testing.T) {
	gi, err := decodeGhostInfo(bytes.NewReader([]byte(testdata)))
	if err != nil {
		t.Fatalf("Unexpected Error Encountered (%v)", err)
	}

	if !reflect.DeepEqual(gi, expected) {
		t.Errorf("got %v,\nwant %v", gi, expected)
	}
}

var expected = ghostInfo{
	m: meta{
		ExportedOn: 1388805572000,
		Version:    "003",
	},
	users: []user{
		{
			ID:   json.RawMessage("1"),
			Name: "user's name",
		},
	},
	tags: []tag{
		{
			ID:   json.RawMessage("3"),
			Name: "Colorado Ho!",
		},
		{
			ID:   json.RawMessage("4"),
			Name: "blue",
		},
	},
	posttags: []posttag{
		{
			PostID: json.RawMessage("5"),
			TagID:  json.RawMessage("3"),
		},
		{
			PostID: json.RawMessage("2"),
			TagID:  json.RawMessage("3"),
		},
		{
			PostID: json.RawMessage("24"),
			TagID:  json.RawMessage("4"),
		},
	},
}

var testdata = `
{
	"meta":{ "exported_on":  1388805572000, "version": "003" },
	"data":{
		"tags": [
			{
				"id":           3,
				"name":         "Colorado Ho!",
				"slug":         "colorado-ho",
				"description":  ""
			},
			{
				"id":           4,
				"name":         "blue",
				"slug":         "blue",
				"description":  ""
			}
		],
		"posts_tags": [
			{"tag_id":3, "post_id":5},
			{"tag_id":3, "post_id":2},
			{"tag_id":4, "post_id":24}
		],
		"users": [
			{
				"id":           1,
				"name":         "user's name",
				"slug":         "users-name",
				"email":        "user@example.com",
				"image":        null,
				"cover":        null,
				"bio":          null,
				"website":      null,
				"location":     null,
				"accessibility": null,
				"status":       "active",
				"language":     "en_US",
				"meta_title":   null,
				"meta_description": null,
				"last_login":   null,
				"created_at":   1283780649000,
				"created_by":   1,
				"updated_at":   1286958624000,
				"updated_by":   1
			}
		],
	}
}
`
