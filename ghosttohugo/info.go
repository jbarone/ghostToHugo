package ghosttohugo

import (
	"encoding/json"
	"io"
)

type meta struct {
	ExportedOn int64  `json:"exported_on"`
	Version    string `json:"version"`
}

type user struct {
	ID   json.RawMessage `json:"id"`
	Name string          `json:"name"`
}

type tag struct {
	ID   json.RawMessage `json:"id"`
	Name string          `json:"name"`
}

type posttag struct {
	ID        json.RawMessage `json:"id"`
	PostID    json.RawMessage `json:"post_id"`
	TagID     json.RawMessage `json:"tag_id"`
	SortOrder int             `json:"sort_order,omitempty"`
}

type setting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type data struct {
	Users    []user    `json:"users"`
	Tags     []tag     `json:"tags"`
	PostTags []posttag `json:"posts_tags"`
	Settings []setting `json:"settings"`
}

type info struct {
	Meta     meta `json:"meta"`
	Data     data `json:"data"`
	settings map[string]string
}

func (c *Converter) decodeInfo(r io.Reader) error {
	decoder := json.NewDecoder(r)

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if tok == "db" {
			decoder.Token()
			err = decoder.Decode(&c.info)
			if err != nil {
				return err
			}
		}

		c.info.settings = make(map[string]string)
		for _, setting := range c.info.Data.Settings {
			c.info.settings[setting.Key] = setting.Value
		}

	}

	return nil
}
