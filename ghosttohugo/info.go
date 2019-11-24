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

type info struct {
	meta     meta
	users    []user
	tags     []tag
	posttags []posttag
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

		switch tok {
		case "meta":
			err = decoder.Decode(&c.info.meta)
		case "users":
			err = decoder.Decode(&c.info.users)
		case "tags":
			err = decoder.Decode(&c.info.tags)
		case "posts_tags":
			err = decoder.Decode(&c.info.posttags)
		case "settings":
			var settings []setting
			err = decoder.Decode(&settings)
			if err != nil || len(settings) == 0 {
				break // out of switch
			}
			c.info.settings = make(map[string]string)
			for _, setting := range settings {
				c.info.settings[setting.Key] = setting.Value
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}
