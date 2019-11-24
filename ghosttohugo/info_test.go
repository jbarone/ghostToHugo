package ghosttohugo

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestConverter_decodeInfo(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    info
		wantErr bool
	}{
		// Empty Test
		{"empty", "", info{}, false},

		// Meta Tests
		{"empty_meta", `{"meta": {}}`, info{}, false},
		{
			"version_meta",
			`{"meta": {"version": "0.3.0"}}`,
			info{
				meta: meta{Version: "0.3.0"},
			},
			false,
		},
		{
			"exported_on_meta",
			`{"meta": {"exported_on": 1234567890}}`,
			info{
				meta: meta{ExportedOn: 1234567890},
			},
			false,
		},

		// Users tests
		{
			"single_user",
			`{"users": [{"id":1234,"name":"username"}]}`,
			info{
				users: []user{
					user{json.RawMessage("1234"), "username"},
				},
			},
			false,
		},
		{
			"multiple_users",
			`{"users": [
				{"id":1234,"name":"username1"},
				{"id":4321,"name":"username2"}
			]}`,
			info{
				users: []user{
					user{json.RawMessage("1234"), "username1"},
					user{json.RawMessage("4321"), "username2"},
				},
			},
			false,
		},

		// Tag tests
		{
			"single_tag",
			`{"tags": [{"id":1234,"name":"tagname"}]}`,
			info{
				tags: []tag{
					tag{json.RawMessage("1234"), "tagname"},
				},
			},
			false,
		},
		{
			"multiple_tags",
			`{"tags": [
				{"id":1234,"name":"tagname1"},
				{"id":4321,"name":"tagname2"}
			]}`,
			info{
				tags: []tag{
					tag{json.RawMessage("1234"), "tagname1"},
					tag{json.RawMessage("4321"), "tagname2"},
				},
			},
			false,
		},

		// PostTags tests
		{
			"singe_post_tag",
			`{"posts_tags": [{"id": 1234, "post_id": 4321, "tag_id": 5432}]}`,
			info{
				posttags: []posttag{
					posttag{
						ID:     json.RawMessage("1234"),
						PostID: json.RawMessage("4321"),
						TagID:  json.RawMessage("5432"),
					},
				},
			},
			false,
		},
		{
			"multiple_post_tags",
			`{"posts_tags": [
				{"id": 1234, "post_id": 4321, "tag_id": 5432},
				{"id": 2468, "post_id": 5555, "tag_id": 6666}
			]}`,
			info{
				posttags: []posttag{
					posttag{
						ID:     json.RawMessage("1234"),
						PostID: json.RawMessage("4321"),
						TagID:  json.RawMessage("5432"),
					},
					posttag{
						ID:     json.RawMessage("2468"),
						PostID: json.RawMessage("5555"),
						TagID:  json.RawMessage("6666"),
					},
				},
			},
			false,
		},
		// TODO(joshua): add tests for settings
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Converter{}
			r := bytes.NewReader([]byte(tt.input))
			if err := c.decodeInfo(r); (err != nil) != tt.wantErr {
				t.Errorf("Converter.decodeInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(c.info, tt.want) {
				t.Errorf("Converte.decodeInfo() got %v, want %v", c.info, tt.want)
			}
		})
	}
}
