package ghostToHugo

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

type ghostInfo struct {
	m        meta
	users    []user
	tags     []tag
	posttags []posttag
}

func decodeGhostInfo(r io.Reader) (ghostInfo, error) {
	var gi ghostInfo
	var decoder = json.NewDecoder(r)
	var doneCount int

	for doneCount < 4 {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return gi, err
		}

		switch tok {
		case "meta":
			err = decoder.Decode(&gi.m)
			doneCount++
		case "users":
			err = decoder.Decode(&gi.users)
			doneCount++
		case "tags":
			err = decoder.Decode(&gi.tags)
			doneCount++
		case "posts_tags":
			err = decoder.Decode(&gi.posttags)
			doneCount++
		}

		if err != nil {
			return gi, err
		}
	}

	return gi, nil
}
