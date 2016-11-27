package ghostToHugo

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// GhostToHugo handles the imprt of a Ghot blog export and outputting to
// hugo static blog
type GhostToHugo struct {
	location   *time.Location
	dateformat string
}

// Post is a blog post in Ghost
type Post struct {
	ID              int             `json:"id"`
	Title           string          `json:"title"`
	Slug            string          `json:"slug"`
	Content         string          `json:"markdown"`
	Image           string          `json:"image"`
	Page            json.RawMessage `json:"page"`
	Status          string          `json:"status"`
	MetaDescription string          `json:"meta_description"`
	AuthorID        int             `json:"author_id"`
	PublishedAt     json.RawMessage `json:"published_at"`
	CreatedAt       json.RawMessage `json:"created_at"`

	Published time.Time
	Created   time.Time
	IsDraft   bool
	IsPage    bool
	Author    string
	Tags      []string
}

func (p *Post) populate(gi *ghostInfo, gth *GhostToHugo) {
	p.Published = gth.parseTime(p.PublishedAt)
	p.Created = gth.parseTime(p.CreatedAt)
	p.IsDraft = p.Status == "draft"
	p.IsPage = parseBool(p.Page)

	for _, user := range gi.users {
		if user.ID == p.AuthorID {
			p.Author = user.Name
			break
		}
	}

	for _, pt := range gi.posttags {
		if pt.PostID == p.ID {
			for _, t := range gi.tags {
				if t.ID == pt.TagID {
					p.Tags = append(p.Tags, t.Name)
					break
				}
			}
		}
	}
}

func parseBool(rm json.RawMessage) bool {
	var b bool
	if err := json.Unmarshal(rm, &b); err == nil {
		return b
	}

	var i int
	if err := json.Unmarshal(rm, &i); err == nil {
		return i != 0
	}

	return false
}

type meta struct {
	ExportedOn int64  `json:"exported_on"`
	Version    string `json:"version"`
}

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type posttag struct {
	ID        int `json:"id"`
	PostID    int `json:"post_id"`
	TagID     int `json:"tag_id"`
	SortOrder int `json:"sort_order,omitempty"`
}

type ghostInfo struct {
	m        meta
	users    []user
	tags     []tag
	posttags []posttag
}

// WithLocation sets the location used when working with timestamps
func WithLocation(location *time.Location) func(*GhostToHugo) {
	return func(gth *GhostToHugo) {
		gth.location = location
	}
}

// WithDateFormat sets the date format to use for ghost imports
func WithDateFormat(format string) func(*GhostToHugo) {
	return func(gth *GhostToHugo) {
		gth.dateformat = format
	}
}

// NewGhostToHugo returns a new instance of GhostToHugo
func NewGhostToHugo(options ...func(*GhostToHugo)) (*GhostToHugo, error) {
	gth := new(GhostToHugo)

	// set defaults
	gth.dateformat = time.RFC3339
	gth.location = time.Local

	for _, option := range options {
		option(gth)
	}

	return gth, nil
}

func seekTo(d *json.Decoder, token json.Token) error {
	var tok json.Token
	var err error
	for err == nil && tok != token {
		tok, err = d.Token()
	}
	return err
}

func stripToken(decoder *json.Decoder) error {
	_, err := decoder.Token() // read in delim
	return err
}

func decodeUsers(r io.Reader) ([]user, error) {
	var decoder = json.NewDecoder(r)
	var users []user
	err := seekTo(decoder, "users")
	if err != nil {
		return users, err
	}
	err = stripToken(decoder)
	if err != nil {
		return users, err
	}
	for decoder.More() {
		var u user
		err = decoder.Decode(&u)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}
	err = stripToken(decoder)
	return users, err
}

func decodeTags(r io.Reader) ([]tag, error) {
	var decoder = json.NewDecoder(r)
	var tags []tag
	err := seekTo(decoder, "tags")
	if err != nil {
		return tags, err
	}
	err = stripToken(decoder)
	if err != nil {
		return tags, err
	}
	for decoder.More() {
		var t tag
		err = decoder.Decode(&t)
		if err != nil {
			return tags, err
		}
		tags = append(tags, t)
	}
	err = stripToken(decoder)
	return tags, err
}

func decodePostTags(r io.Reader) ([]posttag, error) {
	var decoder = json.NewDecoder(r)
	var posttags []posttag
	err := seekTo(decoder, "posts_tags")
	if err != nil {
		return posttags, err
	}
	err = stripToken(decoder)
	if err != nil {
		return posttags, err
	}
	for decoder.More() {
		var t posttag
		err = decoder.Decode(&t)
		if err != nil {
			fmt.Println(err)
			return posttags, err
		}
		posttags = append(posttags, t)
	}
	err = stripToken(decoder)
	return posttags, err
}

func (gth *GhostToHugo) importGhost(r io.ReadSeeker) (<-chan Post, error) {

	var decoder = json.NewDecoder(r)
	var err error
	var gi ghostInfo

	err = seekTo(decoder, "meta")
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(&gi.m)
	if err != nil {
		return nil, err
	}

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	gi.users, err = decodeUsers(r)
	if err != nil {
		return nil, err
	}

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	gi.tags, err = decodeTags(r)
	if err != nil {
		return nil, err
	}

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	gi.posttags, err = decodePostTags(r)
	if err != nil {
		return nil, err
	}

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	decoder = json.NewDecoder(r)
	err = seekTo(decoder, "posts")
	if err != nil {
		return nil, err
	}
	err = stripToken(decoder)
	if err != nil {
		return nil, err
	}

	posts := make(chan Post)
	go func(decoder *json.Decoder, posts chan Post) {
		for decoder.More() {
			var p Post
			err = decoder.Decode(&p)
			if err != nil {
				break
			}
			posts <- p
		}
		close(posts)
	}(decoder, posts)

	return posts, nil
}

func (gth *GhostToHugo) parseTime(raw json.RawMessage) time.Time {
	var pt int64
	if err := json.Unmarshal(raw, &pt); err == nil {
		return time.Unix(0, pt*int64(time.Millisecond)).In(gth.location)
	}

	var ps string
	if err := json.Unmarshal(raw, &ps); err == nil {
		t, err := time.ParseInLocation(gth.dateformat, ps, gth.location)
		if err != nil {
			return time.Time{}
		}
		return t
	}

	return time.Time{}
}
