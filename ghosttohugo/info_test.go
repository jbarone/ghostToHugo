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
		{
			"empty_meta",
			`{"db":[{"meta": {}}]}`,
			info{settings: make(map[string]string)},
			false,
		},
		{
			"version_meta",
			`{"db":[{"meta": {"version": "0.3.0"}}]}`,
			info{
				Meta:     meta{Version: "0.3.0"},
				settings: make(map[string]string),
			},
			false,
		},
		{
			"exported_on_meta",
			`{"db":[{"meta": {"exported_on": 1234567890}}]}`,
			info{
				Meta:     meta{ExportedOn: 1234567890},
				settings: make(map[string]string),
			},
			false,
		},

		// Users tests
		{
			"single_user",
			`{"db":[{"data":{"users": [{"id":1234,"name":"username"}]}}]}`,
			info{
				Data: data{
					Users: []user{
						user{json.RawMessage("1234"), "username"},
					},
				},
				settings: make(map[string]string),
			},
			false,
		},
		{
			"multiple_users",
			`{"db":[{"data":{"users": [
				{"id":1234,"name":"username1"},
				{"id":4321,"name":"username2"}
			]}}]}`,
			info{
				Data: data{
					Users: []user{
						user{json.RawMessage("1234"), "username1"},
						user{json.RawMessage("4321"), "username2"},
					},
				},
				settings: make(map[string]string),
			},
			false,
		},

		// Tag tests
		{
			"single_tag",
			`{"db":[{"data":{"tags": [{"id":1234,"name":"tagname"}]}}]}`,
			info{
				Data: data{
					Tags: []tag{
						tag{json.RawMessage("1234"), "tagname"},
					},
				},
				settings: make(map[string]string),
			},
			false,
		},
		{
			"multiple_tags",
			`{"db":[{"data":{"tags": [
				{"id":1234,"name":"tagname1"},
				{"id":4321,"name":"tagname2"}
			]}}]}`,
			info{
				Data: data{
					Tags: []tag{
						tag{json.RawMessage("1234"), "tagname1"},
						tag{json.RawMessage("4321"), "tagname2"},
					},
				},
				settings: make(map[string]string),
			},
			false,
		},

		// PostTags tests
		{
			"singe_post_tag",
			`{"db":[{"data":{
				"posts_tags": [{"id": 1234, "post_id": 4321, "tag_id": 5432}]
			}}]}`,
			info{
				Data: data{
					PostTags: []posttag{
						posttag{
							ID:     json.RawMessage("1234"),
							PostID: json.RawMessage("4321"),
							TagID:  json.RawMessage("5432"),
						},
					},
				},
				settings: make(map[string]string),
			},
			false,
		},
		{
			"multiple_post_tags",
			`{"db":[{"data":{"posts_tags": [
				{"id": 1234, "post_id": 4321, "tag_id": 5432},
				{"id": 2468, "post_id": 5555, "tag_id": 6666}
			]}}]}`,
			info{
				Data: data{
					PostTags: []posttag{
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
				settings: make(map[string]string),
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
